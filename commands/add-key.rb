name = next_arg("Please specify the key name")
if name !~ /[a-zA-Z0-9]+/
  $stderr.puts "Not a valid key name (letters and numbers only)"
  exit 1
end

key = $stdin.read.strip
if key !~ /(ssh-rsa|ssh-dsa) ([a-zA-Z0-9\+\/]+[=]*)/
  $stderr.puts "Not a valid key"
  exit 1
end

type = $1
value = $2

dirname = File.join(find_account_dir(@username), "keys")
FileUtils.mkdir_p(dirname)
filename = File.join(dirname, name)
File.open(filename, "w") do |file|
  file << type << " " << value    
end

sync_keys
