keys = [] 
dirname = File.join(find_account_dir(@username), "keys")
Dir.entries(dirname).each do |name|
  next if name == '.' || name == '..'
  keys << name
end

if OPTIONS[:format] == :yaml
  puts YAML::dump(keys)
else
  keys.each { |k| puts k }
end
