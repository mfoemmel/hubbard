if @username != 'admin'
    $stderr.puts "You don't have permission to do that"
    exit 3
end

Dir.entries(Hubbard::ACCOUNTS_PATH)).each do |account|
  next if account == '.' || account == '..'
  puts account
end
