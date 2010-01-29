users = [] 
Dir.entries(File.join(Hubbard::HUB_DATA, "accounts")).each do |entry|
  next if entry == '.' || entry == '..'
  users << entry
end

if OPTIONS[:format] == :yaml
  puts YAML::dump(users)
else
  users.each { |u| puts u }
end
