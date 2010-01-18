dirname = File.join(find_account_dir(@username), "keys")
Dir.entries(dirname).each do |name|
  next if name == '.' || name == '..'
  puts name
end
