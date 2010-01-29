Dir.entries(File.join(Hubbard::HUB_DATA, "accounts")).each do |entry|
  next if entry == '.' || entry == '..'
  puts entry
end
