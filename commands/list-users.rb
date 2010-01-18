Dir.entries(File.join(HUB_DATA, "projects")).each do |entry|
  next if entry == '.' || entry == '..'
  puts entry
end
