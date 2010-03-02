if @username != 'admin'
  $stderr.puts "You don't have permission to do that"
  exit 3
end

accounts = []
Dir.entries(Hubbard::ACCOUNTS_PATH).each do |account|
  next if account == '.' || account == '..'
  accounts << account
end

accounts.sort!

if OPTIONS[:format] == :yaml
  puts YAML::dump(accounts)
else
  accounts.each { |a| puts a }
end
