param (
    [string]$hostName = 'localhost',
    [int]$port = 8080,
    [string]$message = 'Hello, server!'
)

# Create a new TCP client
$client = New-Object System.Net.Sockets.TcpClient($hostName, $port)

# Get the network stream
$stream = $client.GetStream()

# Write to the stream
$writer = New-Object System.IO.StreamWriter($stream)
$writer.WriteLine($message)
$writer.Flush()

# Read the response
$reader = New-Object System.IO.StreamReader($stream)
$response = $reader.ReadLine()

# Close the connection
$client.Close()

# Output the response
$response