name = next_arg("Please specify the key name")
if name !~ Hubbard::KEY_NAME_REGEX
  error 1, "Not a valid key name (letters and numbers only)"
end

key = $stdin.read.strip
if key !~ Hubbard::KEY_REGEX
  error 1, "Not a valid key"
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
