#!/usr/bin/env ruby
require 'json'

class String
	def black; "\e[30m#{self}\e[0m" end
	def red; "\e[31m#{self}\e[0m" end
	def green; "\e[32m#{self}\e[0m" end
	def brown; "\e[33m#{self}\e[0m" end
	def blue; "\e[34m#{self}\e[0m" end
	def magenta; "\e[35m#{self}\e[0m" end
	def cyan; "\e[36m#{self}\e[0m" end
	def gray; "\e[37m#{self}\e[0m" end
end

conf_path = ARGV[0];

path = conf_path.split('.')
config_file = File.read("conf.jsonw")
config = JSON.parse(config_file)

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

class Label
	@variable
	attr_accessor :variable
end

class DeepWildcard < Label
	@next_path
	attr_accessor :next_path

	def initialize(next_path)
		@next_path= next_path
	end

	def consume(path)
		both(path)[0]
	end

	def both(path)
		return [[], path] if @next_path.nil?
		drop_index = path.each_index.select{|i| path [i] == @next_path}.max
		return [nil, nil] unless drop_index
		return [path.drop(drop_index), path.take(drop_index)]
	end

	def to_s
		"DeepWildcard #{@next_path}"
	end
end

class DirectHit < Label
	@element
	attr_accessor :element
	def initialize(element)
		@element= element
	end

	def consume(path)
		both(path)[0]
	end

	def both(path)
		return [path.drop(1),path.take(1)] if(path[0] == @element)	
		#Let's be explicit about return nil, since we're using it as a poor man's option
		return [nil,nil]
	end

	def to_s
		"DirectHit #{@element}"
	end
end

class EnumHit < Label
	@entries
	attr_accessor :entries
	def initialize(entries)
		@entries = entries
	end

	def consume(path)
		both(path)[0]
	end

	def both(path)
		return [path.drop(1),path.take(1)] if @entries.index(path[0])
		#Let's be explicit about return nil, since we're using it as a poor man's option
		return [nil,nil]
	end

	def to_s
		"EnumHit #{@entries}"
	end
end

class WildcardHit < Label
	def consume(path)
		both(path)[0]
	end

	def both(path)
		return [path.drop(1),path.take(1)] if path[0]
		#Let's be explicit about return nil, since we're using it as a poor man's option
		return [nil,nil]
	end

	def to_s
		"Wildcard Hit"
	end
end

class RegexHit < Label
	@regex
	attr_accessor :regex
	def initialize(regex)
		@regex = /#{Regexp.new(regex)}/
	end

	def consume(path)
		both(path)[0]
	end

	def both(path)
		return [path.drop(1), path.take(1)] if path[0] and path[0].scan(@regex) != []
		#Let's be explicit about return nil, since we're using it as a poor man's option
		return [nil,nil]
	end

	def to_s
		"RegexHit #{@regex}"
	end
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
	processors.each do |processor|
		path = processor.consume(path)
		if(path.nil?)
			return false
		end
	end
	return path.size == 0
end

data = [
	["test.a","test.a",true],
	["test.a,b","test.a",true],
	["test.a,b","test.b",true],
	["test.a,b","test.c",false],
	["test./[a-z]/","test.a",true],
	["test./[a-z]+/","test.a",true],
	["test./[a-z]/","test.A",false],
	["*.a","test.a",true],
	["*.b","test.a",false],
	["no-key","test.a",false],
	["**.a","test.b",false],
	["**.a","test.a",true],
	["**.a.b.c","test.a.b.c",true],
	["a.**.a","a.test.a",true],
	["a.**.a","a.test.1.a",true],
	["a.**.a","a.test.1.2.a",true],
	["a.**.a","a.a",true],
	["a.**.a","a.test.a.test.a",true],
	["a.**.a","a.test.a.b",false],
	["a.**.a","b.test.a",false],
	["a.**.a","test.a",false],
	["a.**.a","test.a.a",false],
	["**.other","a",false],
]

data.each do |entry| 
	key, path_text, expected = entry
	test_path = path_text.split(".")
	actual = match(key,test_path)
	output =  "'#{key}', '#{path_text}', Expected: '#{expected}' Actual:'#{actual}'  \n"
	output = actual == expected ? "PASS ".green + output : "FAIL ".red + output
	#print output
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

sort_examples = [
	["a","b","c"],
	["a","a.a","a.a.a"],
	["c","b.c","a.b.c","**.a"],
	["a","a,b","/abc/","*","**"],
	["a","*","a.b","a.*","*.b"],
	["a","*","a.b","a.*","*.b"],
	["a","*","a.b","a.*","*.b","*.c"],
]

sort_examples.each do |entry|
	shuffled = entry.shuffle
	sorted = shuffled.sort {|a,b| compare_patterns a, b}
	output = "Expected: '#{entry}' Actual: '#{sorted}' Shuffled: '#{shuffled}'\n"
	output = sorted == entry ? "PASS ".green + output : "FAIL ".red + output
	#print output
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
			variables = {}
			processors.each do |processor|
				local_path,accum = processor.both(local_path)
				variables[processor.variable] = accum.join(".") if processor.variable
				#print "#{local_path},#{accum}\n"
			end
			#print "#{sub_path}, #{hit_key}, #{variables}\n"
			variables.each do |var, var_val|
				next_value = next_value.gsub("${#{var}}",var_val)
			end
			return next_value 
		end
		next_result = recursive_search(path.drop(sub_path.size),next_value)
		return next_result unless next_result.nil?
	end
	return nil
end

result = recursive_search(path,config)
if result.nil?
	puts "Could not find #{conf_path}"
	exit 1
end
puts result
