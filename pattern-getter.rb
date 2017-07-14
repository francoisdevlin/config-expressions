#!/usr/bin/env ruby
require 'json'
require 'optparse'
require_relative 'matchers'
require_relative 'functions'

command = ARGV[0]
args = ARGV.drop(1)

config_file_path="conf.jsonw"
opt_parser = OptionParser.new do |opts|
      opts.banner = "Usage: example.rb [options]"

      opts.separator ""
      opts.separator "Specific options:"

      # Mandatory argument.
      opts.on("-f", "--file FILE",
              "Reads json config from file, defaults to #{config_file_path}") do |conf_file|
      	config_file_path = conf_file 
      end
end

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

opt_parser.parse!(args)
conf_path = args[0];

path = conf_path.split('.')
config_file = File.read(config_file_path)
config = JSON.parse(config_file)

start_state = PatternState.new
start_state.path = path

result = recursion_2(start_state,config)


def explain(state)
	rule = state.evaluated_path.join(".")
	if state.state == :complete
		return "HIT ".green + ": #{rule} VALUE: '#{state.value}'"
	else
		return "MISS".red   + ": #{rule}, ignoring any children".gray
	end
end

if command == "lookup" 
	result = result.reject{|a,b| b.state == :missing }
	if result.nil? or result.size == 0 or result[0][1].state != :complete
		puts "Could not find #{conf_path}"
		exit 1
	end
	puts result[0][1].value
	exit 0
elsif command == "explain"
	i = 0
	print "The following rules were evaluated in this order, the first hit is returned\n"
	final_result = nil
	result.each do |pattern,state|
		i+= 1

		print "Rule %5d: #{explain state}\n" %[i]
		final_result = state.value if state.state == :complete and final_result.nil?
	end
	if final_result.nil?
		print "Could not find #{conf_path}"
	else
		puts final_result
	end
	exit 0
end
