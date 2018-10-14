import socket
import json

PORT = 16379


def test_command(command, key="", dType="", value=None):
	set_protocol = {
		'Command': command,
		'Data': {
		'Key': key,
		'Data': {
			'DType': dType,
			'Value': value
			}
		}
	}
	client_message = json.dumps(set_protocol)
	print("Message from client: ",client_message)
	# create a client socket
	client_sock = socket.socket(family=socket.AF_INET, type=socket.SOCK_STREAM, proto=0)
	# connect to server listening on port 8000
	client_sock.connect(('127.0.0.1', PORT))
	# send (request) data to server
	client_sock.send(client_message.encode())
	# receive data from server
	server_message = client_sock.recv(1024)
	# print received message
	print("Message from server: ", server_message.decode())
	# close the client socket
	client_sock.close()


if __name__ == "__main__":
	test_command("GET", '1')
	test_command('SET', '1','string', '111110000111100001111')
	test_command("GET", '1')
	test_command('SET', '2', 'string', '222220000222200002222')
	test_command("GET", '2')
	test_command('DEL', '1')
	test_command("GET", '1')
	test_command("EVICT")
	test_command("GET", '2')
