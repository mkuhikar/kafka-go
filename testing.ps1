# # Connect to TCP server
# $client = New-Object System.Net.Sockets.TcpClient
# $client.Connect("localhost", 9092)
# $stream = $client.GetStream()

# # Send a fake request to the server
# $request = [System.Text.Encoding]::ASCII.GetBytes("Placeholder request")
# $stream.Write($request, 0, $request.Length)

# # Read 8 bytes (Kafka response)
# $response = New-Object byte[] 8
# $bytesRead = $stream.Read($response, 0, 8)

# # Print the response in hex format
# Write-Host "Received $bytesRead bytes:"
# $response | ForEach-Object { "{0:X2}" -f $_ } | ForEach-Object -Begin { $i = 0 } -Process {
#     if ($i % 16 -eq 0) { Write-Host "`n" -NoNewline }
#     Write-Host "$_ " -NoNewline
#     $i++
# }
# Write-Host ""

# # Cleanup
# $stream.Close()
# $client.Close()

# Convert the hex string to a byte array
$hexString = "00000023001200046f7fc66100096b61666b612d636c69000a6b61666b612d636c6904302e3100"
$bytes = for ($i = 0; $i -lt $hexString.Length; $i += 2) {
    [Convert]::ToByte($hexString.Substring($i, 2), 16)
}

# Connect to TCP server
$client = New-Object System.Net.Sockets.TcpClient
$client.Connect("localhost", 9092)
$stream = $client.GetStream()

# Send the decoded Kafka binary request
$stream.Write($bytes, 0, $bytes.Length)

# Read 8-byte response
$response = New-Object byte[] 8
$bytesRead = $stream.Read($response, 0, 8)

# Print the response in hex format
Write-Host "Received $bytesRead bytes:"
$response | ForEach-Object { "{0:X2}" -f $_ } | ForEach-Object -Begin { $i = 0 } -Process {
    if ($i % 16 -eq 0) { Write-Host "`n" -NoNewline }
    Write-Host "$_ " -NoNewline
    $i++
}
Write-Host ""

# Cleanup
$stream.Close()
$client.Close()
