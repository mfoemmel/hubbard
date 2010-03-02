if @username != 'admin'
  error 3, "You don't have permission to do that"
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
