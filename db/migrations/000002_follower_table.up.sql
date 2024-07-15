CREATE TABLE followers(
                          follower_id UUID REFERENCES users(id),
                          following_id UUID REFERENCES users(id),
                          followed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);