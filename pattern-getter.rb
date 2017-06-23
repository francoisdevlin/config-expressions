#!/usr/bin/env ruby
require 'json'

conf_path = ARGV[0];

path = conf_path.split('.')
config_file = File.read("conf.jsonw")
config = JSON.parse(config_file)

def has_element(key,local_map,globals)
	return true if local_map.key? key
	return false unless local_map.key? '!include'
	local_map['!include'].each do |reference|
		referenced_map = globals["!" + reference]
		return true if referenced_map.key? key
	end
	return false
end

def get_element(key,local_map,globals)
	return local_map[key] if local_map.key? key
	return nil unless local_map.key? '!include'
	local_map['!include'].each do |reference|
		referenced_map = globals["!" + reference]
		return referenced_map[key] if referenced_map.key? key
	end
	return nil
end

def build_the_guys(path,config)
	visited_hash = config
	nested_hash = []
	nested_hash << visited_hash
	i=0;
	built_string=""
	while i<path.size  do
		built_string+=path[i]
		if has_element(built_string,visited_hash, config)
			visited_hash = get_element(built_string, visited_hash, config)
			nested_hash << visited_hash
			built_string=""
		else
			built_string+='.'
		end
		i+=1
	end	
	return built_string,nested_hash
end

def direct_hit(path,config)
	built_string,nested_hash = build_the_guys(path,config)
	if built_string == ""
		return true, nested_hash.last
	end
	return false, nil
end

def wildcard_hit(path,config)
	built_string, nested_hash = build_the_guys(path,config)
	nested_hash.reverse.each do |stack_entry|
		stack_entry = stack_entry.select do | key, value |
			key.start_with? '*.'
		end
		stack_entry = stack_entry.select do | key, value |
			clean_key = key.sub '*.', ''
			key_path = clean_key.split(".")	
			reverse_path = path.reverse
			reverse_key_path = key_path.reverse
			reverse_path[0] == reverse_key_path[0]
		end
		#TODO: SORT TO PREFER LONGER EXACT MATCHES
		if stack_entry.size == 1
			return true, stack_entry.first[1]
		end
	end
	return false, nil
end

def ends_with(path,config)
	built_string, nested_hash = build_the_guys(path,config)
	#We know that built string won't be a match...
	nested_hash.reverse.each do |stack_entry|
		stack_entry = stack_entry.select do | key, value |
			key.start_with? '**'
		end
		stack_entry = stack_entry.select do | key, value |
			clean_key = key.sub '**.', ''
			key_path = clean_key.split(".")	
			reverse_path = path.reverse
			reverse_key_path = key_path.reverse
			#TODO: Make this much much smarter
			reverse_path[0] == reverse_key_path[0]
		end
		#TODO: SORT TO PREFER LONGER EXACT MATCHES
		if stack_entry.size == 1
			return true, stack_entry.first[1]
		end
	end
	return false, nil
end

attempts = [
	[:direct_hit, Proc.new{ |p,c| direct_hit(p,c)}],
	[:wildcard_hit, Proc.new{ |p,c| wildcard_hit(p,c)}],
	[:ends_with, Proc.new{ |p,c| ends_with(p,c)}],
]

attempts.each do |data|
	name,attempt = data
	success, value = attempt.call(path,config)
	if success
		puts value
		exit 0;
	end
end

puts "Could not find key"
exit 1
