package main

import (
    "container/heap"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "sort"
    "sync"
    "time"
)

type User struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    PostCount int    `json:"post_count"`
}

type UserHeap []User

func (h UserHeap) Len() int           { return len(h) }
func (h UserHeap) Less(i, j int) bool { return h[i].PostCount > h[j].PostCount }
func (h UserHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *UserHeap) Push(x interface{}) { *h = append(*h, x.(User)) }
func (h *UserHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

type Post struct {
    ID           string    `json:"id"`
    UserID       string    `json:"userid"`
    Content      string    `json:"content"`
    CommentCount int       `json:"comment_count"`
    Timestamp    time.Time `json:"timestamp"`
}

type PostHeap []Post

func (h PostHeap) Len() int           { return len(h) }
func (h PostHeap) Less(i, j int) bool { return h[i].CommentCount > h[j].CommentCount }
func (h PostHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *PostHeap) Push(x interface{}) { *h = append(*h, x.(Post)) }
func (h *PostHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

type DataStore struct {
    users     map[string]User
    topUsers  UserHeap
    posts     map[string]Post
    topPosts  PostHeap
    latest    []Post
    mu        sync.Mutex
}

func NewDataStore() *DataStore {
    return &DataStore{
        users:    make(map[string]User),
        topUsers: make(UserHeap, 0),
        posts:    make(map[string]Post),
        topPosts: make(PostHeap, 0),
        latest:   make([]Post, 0),
    }
}

func (ds *DataStore) UpdateUser(user User) {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    if _, exists := ds.users[user.ID]; exists {
        for i, u := range ds.topUsers {
            if u.ID == user.ID {
                ds.topUsers[i] = user
                heap.Fix(&ds.topUsers, i)
                break
            }
        }
    } else {
        heap.Push(&ds.topUsers, user)
        if ds.topUsers.Len() > 5 {
            heap.Pop(&ds.topUsers)
        }
    }
    ds.users[user.ID] = user
}

func (ds *DataStore) UpdatePost(post Post) {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    if _, exists := ds.posts[post.ID]; exists {
        for i, p := range ds.topPosts {
            if p.ID == post.ID {
                ds.topPosts[i] = post
                heap.Fix(&ds.topPosts, i)
                break
            }
        }
        for i, p := range ds.latest {
            if p.ID == post.ID {
                ds.latest[i] = post
                sort.Slice(ds.latest, func(i, j int) bool {
                    return ds.latest[i].Timestamp.After(ds.latest[j].Timestamp)
                })
                break
            }
        }
    } else {
        heap.Push(&ds.topPosts, post)
        ds.latest = append(ds.latest, post)
        sort.Slice(ds.latest, func(i, j int) bool {
            return ds.latest[i].Timestamp.After(ds.latest[j].Timestamp)
        })
        if len(ds.latest) > 5 {
            ds.latest = ds.latest[:5]
        }
    }
    ds.posts[post.ID] = post
}

func (ds *DataStore) FetchDataFromTestServer() {
    users := map[string]string{
        "1": "John Doe",
        "2": "Jane Doe",
        "3": "Alice Smith",
        "4": "Bob Johnson",
        "5": "Charlie Brown",
    }
    for userID, name := range users {
        user := User{ID: userID, Name: name, PostCount: rand.Intn(10) + 1}
        ds.UpdateUser(user)
        for i := 0; i < user.PostCount; i++ {
            postID := fmt.Sprintf("%s%d", userID, i)
            post := Post{
                ID:           postID,
                UserID:       userID,
                Content:      fmt.Sprintf("Post %d by %s", i, name),
                CommentCount: rand.Intn(5),
                Timestamp:    time.Now().Add(-time.Duration(rand.Intn(3600)) * time.Second),
            }
            ds.UpdatePost(post)
        }
    }
    log.Println("Simulated data loaded")
}

func (ds *DataStore) TopUsersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    ds.mu.Lock()
    defer ds.mu.Unlock()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ds.topUsers)
}

func (ds *DataStore) PostsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    ds.mu.Lock()
    defer ds.mu.Unlock()
    w.Header().Set("Content-Type", "application/json")
    queryType := r.URL.Query().Get("type")

    switch queryType {
    case "popular":
        if ds.topPosts.Len() == 0 {
            json.NewEncoder(w).Encode([]Post{})
            return
        }
        maxCount := ds.topPosts[0].CommentCount
        var topPosts []Post
        for _, p := range ds.posts {
            if p.CommentCount == maxCount {
                topPosts = append(topPosts, p)
            }
        }
        json.NewEncoder(w).Encode(topPosts)
    case "latest":
        json.NewEncoder(w).Encode(ds.latest)
    default:
        http.Error(w, "Invalid type parameter. Use 'popular' or 'latest'", http.StatusBadRequest)
    }
}

func main() {
    ds := NewDataStore()
    ds.FetchDataFromTestServer()
    go func() {
        for {
            time.Sleep(30 * time.Second)
            ds.FetchDataFromTestServer()
        }
    }()
    http.HandleFunc("/users", ds.TopUsersHandler)
    http.HandleFunc("/posts", ds.PostsHandler)
    log.Println("Server starting on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}