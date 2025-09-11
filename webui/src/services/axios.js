import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

// Add a request interceptor to include authentication headers
instance.interceptors.request.use(
	(config) => {
		// Get user ID from localStorage
		const userId = localStorage.getItem('userId');
		
		// Add Authorization header if user is logged in
		// Skip for public endpoints like /session and /liveness
		const publicEndpoints = ['/session', '/liveness'];
		const isPublicEndpoint = publicEndpoints.some(endpoint => 
			config.url.includes(endpoint)
		);
		
		if (userId && !isPublicEndpoint) {
			config.headers.Authorization = `Bearer ${userId}`;
		}
		
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

// Add a response interceptor to handle authentication errors
instance.interceptors.response.use(
	(response) => {
		return response;
	},
	(error) => {
		// If we get a 401 Unauthorized, clear localStorage and reload
		if (error.response && error.response.status === 401) {
			localStorage.removeItem('userId');
			localStorage.removeItem('currentUsername');
			// Optionally redirect to login or reload the page
			window.location.reload();
		}
		return Promise.reject(error);
	}
);

export default instance;
