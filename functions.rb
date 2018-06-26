require_relative 'matchers'

def get_label_info(key)
	matches = key.scan(/(.*)\$(\w+)$/)
	return key,nil if matches == []
	return matches[0]	
end

def parse_processors(key)
	split_keys = key.split('.')
	processors = split_keys.each_index.collect do |index|
		label, variable= get_label_info(split_keys[index])
		result = nil
		if label == "**"
			next_label = split_keys[index+1]
			result = DeepWildcard.new(next_label)
		elsif label == "*"
			result = WildcardHit.new()
		elsif label.start_with?("/") and label.end_with?("/")
			result = RegexHit.new(label[1..label.size()-2])
		elsif label.include? ","
			result = EnumHit.new(label.split(","))
		else
			result = DirectHit.new(label)
		end
		result.variable = variable
		result.label=split_keys[index]
		result
	end
	return processors
end

def match_state(key,state)
	processors = parse_processors(key)
	processors.each do |processor|
		return state if state.path.size == 0
		state = processor.next(state)
		if(state.path.nil?)
			state.state = :missing
			return state
		end
	end
	state.state = :complete if state.path.size == 0
	return state
end

def determine_match_states(start_state,config)
	sorted_matches = []
	if(config.instance_of? Hash)
		sorted_matches = config.keys.sort {|a,b| compare_patterns a, b}
	elsif(config.instance_of? Array)
		sorted_matches = (0..(config.size-1)).map {|x| x.to_s}
	else
	end
	previous_iteration_pattern = ""
	previous_iteration_state = nil
	results = sorted_matches.collect do | match |
		same_class = compare_patterns(previous_iteration_pattern,match) == 0
		next_state = match_state(match,start_state)
		if 1 && 
			same_class && 
			!previous_iteration_state.nil? && 
			[:complete,:incomplete].include?(next_state.state) &&
			[:complete,:incomplete,:collision].include?(previous_iteration_state.state)
			previous_iteration_state.state = :collision
			next_state.state = :collision
		end
		previous_iteration_pattern = match
		previous_iteration_state = next_state
		[match, next_state]
	end
	return results
end


def compare_patterns(left,right)
	deep_wildcard_penalty=10000
	symbol_order = [
		"DirectHit",
		"EnumHit",
		"RegexHit",
		"WildcardHit",
		"DeepWildcard",
	]
	left_processors = parse_processors(left)
	right_processors = parse_processors(right)
	left_size = left_processors.size
	left_processors.each {|processor| left_size+=deep_wildcard_penalty if processor.instance_of? DeepWildcard}
	right_size = right_processors.size
	right_processors.each {|processor| right_size+=deep_wildcard_penalty if processor.instance_of? DeepWildcard}
	result = left_size.<=> right_size
	return result if result !=0
	(0..left_size).each do |index|
		leftProc = left_processors[index]
		rightProc = right_processors[index]
		leftWeight = symbol_order.index leftProc.class.name
		rightWeight = symbol_order.index rightProc.class.name
		result = leftWeight.<=> rightWeight
		return result if result !=0
	end
	(0..left_size).each do |index|
		leftProc = left_processors[index]
		rightProc = right_processors[index]
		if leftProc.class.name == "DirectHit"
			result = leftProc.element.<=>rightProc.element
			return result if result !=0
		end
	end
	return 0
end

def pretty_key(pattern,state)
	output = state.evaluated_path.join(".")
	return [output,pattern].max {|x| x.size}
end

def get_next_value(value,pattern)
	return value[pattern] if value.instance_of? Hash
	return value[pattern.to_i] if value.instance_of? Array
	nil	
end

def recursion_2(input_state,value)
	results = determine_match_states(input_state,value)
	output = []
	results.each do |pattern,state|
		if state.state == :complete
			next_value = get_next_value(value,pattern)
			if next_value.nil?
				state.state = :key_miss_bro if next_value.nil?
			elsif next_value.instance_of? Array
				state.state = :incomplete if next_value.instance_of? Array
			elsif next_value.instance_of? Hash
				state.state = :incomplete if next_value.instance_of? Hash
			end
			state.variables.each do |var, var_val|
				next_value = next_value.gsub("${#{var}}",var_val)
			end
			state.value = next_value
			output << [pretty_key(pattern,state), state]
		elsif state.state == :incomplete && state.path.size > 0
			next_value = get_next_value(value,pattern)
			if next_value.nil?
				state.state = :key_miss_bro if next_value.nil?
				output << [pretty_key(pattern,state), state]
				next
			end
			state.locality << pattern
			output.concat recursion_2(state,next_value)
		else 
			output << [pretty_key(pattern,state), state]
		end
	end
	return output
end

#TODO - Revive this
def hash_lambda_factory(curried_state,conf)
	return Proc.new  do |h, k|
		next_state = PatternState.new
		next_state.evaluated_path = curried_state.evaluated_path.clone
		next_state.variables = curried_state.variables.clone
		next_state.locality = curried_state.locality.clone
		next_state.path = curried_state.path.clone
		next_state.path << k
		result = recursion_2(next_state,conf).reject{|a,b| b.state == :missing}
		result.each {|a,b| print "Result:#{b}\n"}

		if result.size == 0
			next
		end
		highest_expression = result[0][0]
		highest_result = result[0][1]

		if highest_result.state == :missing
			#h[k] = nil
		elsif highest_result.state == :collision
			raise "Collision for '#{k}'"
		elsif highest_result.state == :incomplete
			h[k] = Hash.new(&hash_lambda_factory(next_state,conf))
		else
			h[k] = highest_result.value
		end
	end
end

def hash_factory(conf)
	next_state = PatternState.new
	next_state.path = []
	return Hash.new &hash_lambda_factory(next_state,conf)
end
