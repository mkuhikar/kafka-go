# # # Connect to TCP server
# # $client = New-Object System.Net.Sockets.TcpClient
# # $client.Connect("localhost", 9092)
# # $stream = $client.GetStream()

# # # Send a fake request to the server
# # $request = [System.Text.Encoding]::ASCII.GetBytes("Placeholder request")
# # $stream.Write($request, 0, $request.Length)

# # # Read 8 bytes (Kafka response)
# # $response = New-Object byte[] 8
# # $bytesRead = $stream.Read($response, 0, 8)

# # # Print the response in hex format
# # Write-Host "Received $bytesRead bytes:"
# # $response | ForEach-Object { "{0:X2}" -f $_ } | ForEach-Object -Begin { $i = 0 } -Process {
# #     if ($i % 16 -eq 0) { Write-Host "`n" -NoNewline }
# #     Write-Host "$_ " -NoNewline
# #     $i++
# # }
# # Write-Host ""

# # # Cleanup
# # $stream.Close()
# # $client.Close()

# # Convert the hex string to a byte array
# $hexString = "000000230012674a4f74d28b00096b61666b612d636c69000a6b61666b612d636c6904302e3100"
# $bytes = for ($i = 0; $i -lt $hexString.Length; $i += 2) {
#     [Convert]::ToByte($hexString.Substring($i, 2), 16)
# }

# # Connect to TCP server
# $client = New-Object System.Net.Sockets.TcpClient
# $client.Connect("localhost", 9092)
# $stream = $client.GetStream()

# # Send the decoded Kafka binary request
# $stream.Write($bytes, 0, $bytes.Length)

# # Read 8-byte response
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


# Kafka ApiVersionsRequest (API key 18, version 4), with:
# - Correlation ID: 0x674a4f74 (just an example)
# - Client ID: "kafka-cli"
# - ClientSoftwareName: "kafka-cli"
# - ClientSoftwareVersion: "0.1"
# Hex string was extracted from the Go test logic

$hexString = "000000230012674a4f7400096b61666b612d636c69000a6b61666b612d636c6904302e3100"

# Convert hex string to byte array
$bytes = for ($i = 0; $i -lt $hexString.Length; $i += 2) {
    [Convert]::ToByte($hexString.Substring($i, 2), 16)
}

# Connect to Kafka server on port 9092
$client = New-Object System.Net.Sockets.TcpClient
$client.Connect("localhost", 9092)
$stream = $client.GetStream()

# Send the Kafka ApiVersionsRequest
$stream.Write($bytes, 0, $bytes.Length)

# Read the first 4 bytes to get the response length
$responseSizeBytes = New-Object byte[] 4
$stream.Read($responseSizeBytes, 0, 4) | Out-Null
[Array]::Reverse($responseSizeBytes)
$responseLength = [System.BitConverter]::ToUInt32($responseSizeBytes, 0)
# Read the rest of the response
$responseBody = New-Object byte[] $responseLength
$totalRead = 0
while ($totalRead -lt $responseLength) {
    $read = $stream.Read($responseBody, $totalRead, $responseLength - $totalRead)
    if ($read -eq 0) { break }
    $totalRead += $read
}

# Print the full response in hex format
Write-Host "Received $($totalRead + 4) bytes (including size prefix):"
$responseSizeBytes + $responseBody | ForEach-Object { "{0:X2}" -f $_ } | ForEach-Object -Begin { $i = 0 } -Process {
    if ($i % 16 -eq 0) { Write-Host "`n" -NoNewline }
    Write-Host "$_ " -NoNewline
    $i++
}
Write-Host ""

# Cleanup
$stream.Close()
$client.Close()
