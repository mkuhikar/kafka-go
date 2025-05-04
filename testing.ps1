# Connect to TCP server
$client = New-Object System.Net.Sockets.TcpClient
$client.Connect("localhost", 9092)
$stream = $client.GetStream()

# Send a fake request to the server
$request = [System.Text.Encoding]::ASCII.GetBytes("Placeholder request")
$stream.Write($request, 0, $request.Length)

# Read 8 bytes (Kafka response)
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
