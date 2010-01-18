name = next_arg("Please specify the key name")
if name !~ /[a-zA-Z0-9]+/
  $stderr.puts "Not a valid key name (letters and numbers only)"
  exit 1
end

dirname = File.join(find_account_dir(@username), "keys")
FileUtils.mkdir_p(dirname)

filename = File.join(dirname, name)
if !File.exist?(filename)
  $stderr.puts "Key not found"
  exit 1
end

unless FileUtils.rm(filename)
  $stderr.puts "Unable to delete key"
  exit 1
end

sync_keys
