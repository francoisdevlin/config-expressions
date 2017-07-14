#!/usr/bin/env ruby
require_relative 'matchers'
require_relative 'functions'

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

data = [
	["test","test.a","incomplete"],
	["g","g.a","incomplete"],
	["test.a","test.a","complete"],
	["test.a,b","test.a","complete"],
	["test.a,b","test.b","complete"],
	["test.a,b","test.c","missing"],
	["test./[a-z]/","test.a","complete"],
	["test./[a-z]+/","test.a","complete"],
	["test./[a-z]/","test.A","missing"],
	["*.a","test.a","complete"],
	["*.b","test.a","missing"],
	["no-key","test.a","missing"],
	["**.a","test.b","missing"],
	["**.a","test.a","complete"],
	["**.a.b.c","test.a.b.c","complete"],
	["a.**.a","a.test.a","complete"],
	["a.**.a","a.test.1.a","complete"],
	["a.**.a","a.test.1.2.a","complete"],
	["a.**.a","a.a","complete"],
	["a.**.a","a.test.a.test.a","complete"],
	["a.**.a","a.test.a.b","incomplete"],
	["a.**.a","b.test.a","missing"],
	["a.**.a","test.a","missing"],
	["a.**.a","test.a.a","missing"],
	["**.other","a","missing"],
]

data.each do |entry| 
	key, path_text, expected = entry
	test_path = path_text.split(".")
	actual = match(key,test_path)
	output =  "'#{key}', '#{path_text}', Expected: '#{expected}' Actual:'#{actual}'  \n"
	output = actual.to_s == expected ? "PASS ".green + output : "FAIL ".red + output
	print output
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
	print output
end

conf_dict = {
	"a" => "value_a",
	"b" => "value_b",
	"c,d$enum" => "enum_${enum}",
	"e.*$app.user" => "e_${app}_user",
	"f" => {
		"a" => "value_f_a",
		"/10\\w+/$value" => "10_something_${value}"
	},
	"g.a" => "value_g_a",
	"h,i$enum" => {
		"a" => "nested_${enum}_a",
	},
}
=begin
h = hash_factory(conf_dict)
hash_examples = [
	["a","value_a"],
	["b","value_b"],
	["c","enum_c"],
	["d","enum_d"],
	["e.bacon.user","e_bacon_user"],
	["f.a","value_f_a"],
	["f.10bacon","10_something_10bacon"],
	["g.a","value_g_a"],
]

hash_examples.each do |entry|
	path, expected = entry
	actual = h
	path.split(".").each do	|label|
		actual = actual[label]
	end
	#sorted = shuffled.sort {|a,b| compare_patterns a, b}
	output = "Path: #{path} Expected: '#{expected}' Actual: '#{actual}'\n"
	output = actual == expected ? "PASS ".green + output : "FAIL ".red + output
	#print output
end

start_state = PatternState.new
start_state.path = ["h","a"]
#start_state.path = ["e","a","user"]
#list = determine_match_states(start_state,conf_dict)
list = recursion_2(start_state,conf_dict)
#list = list.reject{|a,b| b.state == :missing }
list.each do |a,b|
	print "#{a}, #{b}\n"
end
=end

#h = hash_factory(conf_dict)
#puts h["a"]
#puts h["b"]
#puts h["c"]
#puts h["d"]
#puts h["g"]["a"]
#puts h["f"]["a"]
#puts h["e"]["bacon"]["user"]
