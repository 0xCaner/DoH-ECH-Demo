import asyncio
import websockets
import ssl

async def echo(websocket, path):
    async for message in websocket:
        print(f"Received message: {message}")
        await websocket.send(f"Echo: {message}")

async def main():
    ssl_context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
    ssl_context.load_cert_chain(certfile='./cert/cert.pem', keyfile='./cert/key.pem')
    server = await websockets.serve(echo, "0.0.0.0", 443, ssl=ssl_context)
    print("WebSocket server started on wss://0.0.0.0:443")
    await server.wait_closed()

if __name__ == "__main__":
    asyncio.run(main())