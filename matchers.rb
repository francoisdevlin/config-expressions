class Label
	@variable
	attr_accessor :variable

	def consume(path)
		both(path)[0]
	end

	def next(state)
		output = PatternState.new
		return output if state.path.nil?
		rest, consumed = both(state.path)
		return output if rest.nil?
		output.path = rest
		output.evaluated_path = state.evaluated_path.clone
		output.evaluated_path << consumed
		output.variables = state.variables.clone
		output.variables[@variable]=consumed.join(".") if @variable
		return output
	end

end

class DeepWildcard < Label
	@next_path
	attr_accessor :next_path

	def initialize(next_path)
		@next_path= next_path
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

	def both(path)
		return [path.drop(1), path.take(1)] if path[0] and path[0].scan(@regex) != []
		#Let's be explicit about return nil, since we're using it as a poor man's option
		return [nil,nil]
	end

	def to_s
		"RegexHit #{@regex}"
	end
end
