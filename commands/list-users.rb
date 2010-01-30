if @username != 'admin'
    $stderr.puts "You don't have permission to do that"
    exit 3
end

Dir.entries(File.join(Hubbard::HUB_DATA, 'accounts')).each do |account|
  next if account == '.' || account == '..'
  puts account
end
