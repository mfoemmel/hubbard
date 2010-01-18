Dir.entries(File.join(HUB_DATA, "accounts")).each do |entry|
  next if entry == '.' || entry == '..'
  puts entry
end
