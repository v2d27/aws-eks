# Real-Time Chat Application - Frontend

This is the React.js + TypeScript frontend for the real-time chat application built with Vite.

## Features
- Real-time messaging with WebSocket
- Clean and responsive UI
- Message history display
- Connection status indicator
- Auto-scroll to new messages
- Keyboard support (Enter to send)
- Environment variable support
- Modern development tooling with Vite

## Technology Stack
- React 18 with TypeScript
- Vite for fast development and building
- Native WebSocket API for real-time communication
- Modern CSS for styling
- Node.js 22.16+ required

## Project Structure
```
src/
├── components/         # Reusable React components
├── hooks/             # Custom React hooks
├── services/          # API and WebSocket services
├── styles/            # CSS and styling files
├── types.ts           # TypeScript type definitions
├── utils/             # Utility functions
├── App.tsx            # Main application component
├── main.tsx           # Application entry point
└── vite-env.d.ts      # Vite environment types
```

## Setup and Running

### Prerequisites
- Node.js 22.16 or later
- yarn

### Development Setup

1. Install dependencies:
```bash
yarn install
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Configure environment variables in `.env`:
```env
VITE_WS_HOST=localhost
VITE_WS_PORT=8080
VITE_WS_PROTOCOL=ws
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_NAME=Chat Application
VITE_DEBUG=true
```

4. Start the development server:
```bash
yarn start
```

The application will run on http://localhost:3000

### Production Build

```bash
# Build for production
yarn build

# Preview production build
yarn preview
```

## Environment Variables

All environment variables must be prefixed with `VITE_`:

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_WS_HOST` | WebSocket host | `localhost` |
| `VITE_WS_PORT` | WebSocket port | `8080` |
| `VITE_WS_PROTOCOL` | WebSocket protocol | `ws` |
| `VITE_API_BASE_URL` | API base URL | `http://localhost:8080` |
| `VITE_APP_NAME` | Application name | `Chat Application` |
| `VITE_DEBUG` | Debug mode | `true` |

## Testing
1. Make sure the backend server is running on port 8080
2. Open the application in two browser tabs
3. Send messages from either tab to see real-time communication

## Contributing

- All contributings are welcome!

## License

- Apache License 2.0, see [LICENSE](./LICENSE).