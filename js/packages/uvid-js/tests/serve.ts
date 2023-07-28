import http from 'node:http';
import handler from 'serve-handler';

export async function serve(dir: string) {
	const server = http.createServer(async (request, response) => {
		// You pass two more arguments for config and middleware
		// More details here: https://github.com/vercel/serve-handler#options
		return handler(request, response, {
			public: dir,
		});
	});

	return new Promise((resolve) => {
		server.listen(4000, () => {
			console.log('Listening at http://localhost:4000');
			resolve(null);
		});
	});
}
