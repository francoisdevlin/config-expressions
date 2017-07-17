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
	sorted_matches = config.keys.sort {|a,b| compare_patterns a, b}
	results = sorted_matches.collect do | match |
		[match, match_state(match,start_state)]
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

def recursion_2(input_state,value)
	results = determine_match_states(input_state,value)
	output = []
	results.each do |pattern,state|
		if state.state == :complete
			next_value = value[pattern]
			if next_value.nil?
				state.state = :key_miss_bro if next_value.nil?
			elsif next_value.instance_of? Hash
				state.state = :too_short_bro if next_value.instance_of Hash
			end
			state.variables.each do |var, var_val|
				next_value = next_value.gsub("${#{var}}",var_val)
			end
			state.value = next_value
			output << [state.evaluated_path.join("."), state]
		elsif state.state == :incomplete
			next_value = value[pattern]
			if next_value.nil?
				state.state = :key_miss_bro if next_value.nil?
				output << [state.evaluated_path.join("."), state]
				next
			end
			state.locality << pattern
			output.concat recursion_2(state,next_value)
		else 
			output << [state.evaluated_path.join("."), state]
		end
	end
	return output
end

def detect_ambiguous_match(input)
	result = input.reject{|a,b| b.state == :missing }
	if result.nil? or result.size == 0 or result[0][1].state != :complete
		return input
	end
	highest_expression = result[0][0]
	highest_result = result[0][1]
	collide = Proc.new {|a,b| compare_patterns(highest_expression,a) == 0 and highest_result.locality == b.locality}
	repeats = result.select &collide
	if repeats.size == 1
		return input
	end
	input.each {|a,b| b.state=:collision if collide.call(a,b)}
	return input
end

#TODO - Revive this
#def hash_lambda_factory(curried_path,conf)
	#return Proc.new  do |h, k|
		#local_path = curried_path.clone << k
		#value = recursive_search(local_path,conf)
		#if value == :missing
			#h[k] = nil
		#elsif value == :incomplete
			#h[k] = Hash.new(&hash_lambda_factory(local_path,conf))
		#else
			#h[k] = value
		#end
	#end
#end

#def hash_factory(conf)
	#return Hash.new &hash_lambda_factory([],conf)
#end
