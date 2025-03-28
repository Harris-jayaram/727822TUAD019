# Social Media Analytics Microservice

## Overview
This microservice provides real-time analytical insights into user activity on a social media platform. It retrieves and processes data from a test server API, delivering insights on top users and posts.

## Features
- Fetches and maintains real-time analytics on social media users and posts.
- Provides two API endpoints:
  - `/users`: Retrieves the top five users with the highest number of posts.
  - `/posts?type={popular|latest}`: Retrieves either the most commented post(s) or the five most recent posts.
- Uses heap-based data structures for efficient storage and retrieval of top users and posts.
- Ensures real-time data updates while optimizing API call costs.

## Technologies Used
- Go (Golang)
- HTTP server
- Heap (Priority Queue) for efficient sorting and retrieval
- JSON encoding/decoding
- Concurrency and synchronization using `sync.Mutex`

## Installation and Setup
### Prerequisites
- Go installed on your system.

### Steps
1. Clone the repository:
   ```sh
   git clone https://github.com/your-repo/social-media-analytics.git
   cd social-media-analytics
   ```
2. Run the microservice:
   ```sh
   go run main.go
   ```
3. The server starts on port `8080` and fetches test data at startup.

## API Endpoints
### 1. Get Top Users
**Endpoint:** `/users`
**Method:** `GET`
**Description:** Returns the top 5 users with the highest post count.
**Response Format:**
```json
[
  { "id": "1", "name": "John Doe", "post_count": 10 },
  { "id": "2", "name": "Jane Doe", "post_count": 9 }
]
```

### 2. Get Popular or Latest Posts
**Endpoint:** `/posts?type={popular|latest}`
**Method:** `GET`
**Query Params:**
- `type=popular`: Returns post(s) with the highest comment count.
- `type=latest`: Returns the 5 most recent posts.

**Response Format:**
```json
[
  { "id": "101", "userid": "1", "content": "Hello World!", "comment_count": 5, "timestamp": "2025-03-28T12:00:00Z" }
]
```

## Implementation Details
- Uses a **heap-based approach** (`UserHeap`, `PostHeap`) to efficiently store and retrieve top users and posts.
- Periodically **fetches data from the test server** to keep analytics up-to-date.
- Ensures **concurrent access** to data structures with a `sync.Mutex`.
- Limits **API call frequency** to reduce costs while maintaining real-time analytics.





