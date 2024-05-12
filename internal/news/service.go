package news

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type NewsService struct {
	rdb *redis.Client
    client *NewsClient
}

func NewNewsService(rdb *redis.Client, client *NewsClient) *NewsService {
	return &NewsService{
        rdb,
        client,
    }
}

type FeedAPIResponse struct {
	Article   News `json:"article"`
	Upvotes   int  `json:"upvotes"`
	Downvotes int  `json:"downvotes"`
}

func (s *NewsService) getSources(ctx context.Context) ([]Source, error) {
    sources, err := s.client.sources("fr", "tech")
    if err != nil {
        return nil, err
    }
    return sources, nil
}

func (s *NewsService) getNews(ctx context.Context) ([]FeedAPIResponse, error) {
	news := []*News{
		{
			ID:    "048566e0-edf5-45bc-8645-9f3f9e948e4d",
			Title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			Desc:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus eu pulvinar tortor. Fusce placerat ligula ac nisi sollicitudin rhoncus. In volutpat convallis purus id blandit. Cras dapibus consequat dignissim. Morbi condimentum felis dolor. Morbi aliquet, nisi eu cursus elementum, mauris erat ultricies leo, vitae imperdiet metus neque a mi. Curabitur in auctor enim.",
			Link:  "https://google.com",
			Img:   "https://leclaireur.fnac.com/wp-content/uploads/2023/10/spiderman7.jpg",
		},
		{
			ID:    "1932b014-d87a-4a84-a358-5c775b9c6bf6",
			Title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			Desc:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus eu pulvinar tortor. Fusce placerat ligula ac nisi sollicitudin rhoncus. In volutpat convallis purus id blandit. Cras dapibus consequat dignissim. Morbi condimentum felis dolor. Morbi aliquet, nisi eu cursus elementum, mauris erat ultricies leo, vitae imperdiet metus neque a mi. Curabitur in auctor enim.",
			Link:  "https://google.com",
			Img:   "https://leclaireur.fnac.com/wp-content/uploads/2023/10/spiderman7.jpg",
		},
		{
			ID:    "f474f870-1775-43f3-b3c0-bafb95903ce2",
			Title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			Desc:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus eu pulvinar tortor. Fusce placerat ligula ac nisi sollicitudin rhoncus. In volutpat convallis purus id blandit. Cras dapibus consequat dignissim. Morbi condimentum felis dolor. Morbi aliquet, nisi eu cursus elementum, mauris erat ultricies leo, vitae imperdiet metus neque a mi. Curabitur in auctor enim.",
			Link:  "https://google.com",
			Img:   "https://leclaireur.fnac.com/wp-content/uploads/2023/10/spiderman7.jpg",
		},
	}

	resp := []FeedAPIResponse{}
	for _, article := range news {
		val, err := s.rdb.Get(ctx, article.ID).Result()
		if err == redis.Nil {
			if err := s.rdb.Set(ctx, article.ID, "0:0", 0).Err(); err != nil {
			}
			resp = append(resp, FeedAPIResponse{
				Article:   *article,
				Upvotes:   0,
				Downvotes: 0,
			})
			continue
		}
		if err != nil {
            return nil, fmt.Errorf("something went wrong")
		}
		votes := strings.Split(val, ":")
		upvotes, err := strconv.Atoi(votes[0])
		if err != nil {
			upvotes = 0
		}
		downvotes, err := strconv.Atoi(votes[1])
		if err != nil {
			downvotes = 0
		}
		resp = append(resp, FeedAPIResponse{
			Article:   *article,
			Upvotes:   upvotes,
			Downvotes: downvotes,
		})

	}
    return resp, nil
}
