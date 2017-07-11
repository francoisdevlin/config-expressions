#!/usr/bin/env ruby
require 'json'
require 'optparse'
require_relative 'matchers'
require_relative 'functions'

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

opt_parser.parse!(ARGV)
conf_path = ARGV[0];

path = conf_path.split('.')
config_file = File.read(config_file_path)
config = JSON.parse(config_file)

result = recursive_search(path,config)
if result.nil?
	puts "Could not find #{conf_path}"
	exit 1
end
puts result
