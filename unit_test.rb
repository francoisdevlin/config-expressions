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
