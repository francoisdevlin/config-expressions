require_relative 'matchers'

def determine_matches(local_path,config)
	matches = []
	sorted_matches = config.keys.sort {|a,b| compare_patterns a, b}
	(0..local_path.size).each do |index|
		sub_path = local_path.take(index+1)
		sorted_matches.each do |key|
			result = match(key,sub_path)
			matches << [sub_path,key] if result
		end
	end
	return matches
end

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
		result
	end
	return processors
end

def match(key,path)
	processors = parse_processors(key)
	state = PatternState.new
	state.path = path
	processors.each do |processor|
		state = processor.next(state)
		if(state.path.nil?)
			return false
		end
	end
	return state.path.size == 0
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

def recursive_search(path,value)
	results = determine_matches(path,value)
	results.each do |result|
		sub_path, hit_key = result
		next_value = value[hit_key]
		next if next_value.nil?
		unless next_value.instance_of? Hash
			processors = parse_processors(hit_key)
			local_path = path
			state = PatternState.new
			state.path= path
			processors.each do |processor|
				state = processor.next(state)
			end
			state.variables.each do |var, var_val|
				next_value = next_value.gsub("${#{var}}",var_val)
			end
			return next_value 
		end
		next_result = recursive_search(path.drop(sub_path.size),next_value)
		return next_result unless next_result.nil?
	end
	return nil
end
