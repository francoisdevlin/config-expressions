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

opt_parser.parse!(args)
conf_path = args[0];

path = conf_path.split('.')
config_file = File.read(config_file_path)
config = JSON.parse(config_file)

start_state = PatternState.new
start_state.path = path

result = recursion_2(start_state,config)

if command == "lookup" 
	result = result.reject{|a,b| b.state == :missing }
	if result.nil? or result.size == 0 or result[0][1].state != :complete
		puts "Could not find #{conf_path}"
		exit 1
	end
	puts result[0][1].value
	exit 0
end
