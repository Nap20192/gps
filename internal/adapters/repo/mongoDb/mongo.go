package mongoDb

import (
	"context"
	"fmt"
	"gps/internal/domain/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	client    *mongo.Client
	db        *mongo.Database
	routeColl *mongo.Collection
	usersColl *mongo.Collection
	ctx       context.Context
}

func NewRepository(ctx context.Context, client *mongo.Client, dbName string) (*Repository, error) {

	db := client.Database(dbName)

	coll := db.Collection("routes")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"route_id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, err
	}
	return &Repository{
		client:    client,
		db:        db,
		routeColl: coll,
		usersColl: db.Collection("users"),
		ctx:       ctx,
	}, nil
}

func (m *Repository) CreateRoute(route models.Route) (*mongo.InsertOneResult, error) {
	return m.routeColl.InsertOne(m.ctx, route)
}

func (m *Repository) GetRouteByID(routeID uuid.UUID) (*models.Route, error) {
	var route models.Route
	err := m.routeColl.FindOne(m.ctx, bson.M{"route_id": routeID}).Decode(&route)
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (m *Repository) GetGPSDataLastNSeconds(n int) ([]models.GPSData, error) {
	since := time.Now().Add(-time.Duration(n) * time.Second)
	cursor, err := m.routeColl.Find(m.ctx, bson.M{
		"path.timestamp": bson.M{"$gte": since},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(m.ctx)
	var result []models.GPSData

	for cursor.Next(m.ctx) {
		var route models.Route
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}

		for _, gps := range route.Path {
			if gps.Timestamp.After(since) || gps.Timestamp.Equal(since) {
				result = append(result, gps)
			}
		}
	}

	return result, nil
}

func (m *Repository) AddGPSDataToRoute(routeID uuid.UUID, gps models.GPSData) (*mongo.UpdateResult, error) {
	update := bson.M{
		"$push": bson.M{"path": gps},
	}
	return m.routeColl.UpdateOne(m.ctx, bson.M{"route_id": routeID}, update)
}

func (m *Repository) Close() error {
	return m.client.Disconnect(m.ctx)
}

func (m *Repository) CreateUser(username, email, passwordHash string) (uuid.UUID, error) {
	id := uuid.New()
	user := models.User{
		UserID:       id,
		Username:     username,
		PasswordHash: passwordHash,
	}
	_, err := m.usersColl.InsertOne(m.ctx, user)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (m *Repository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := m.usersColl.FindOne(m.ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	if user.UserID == uuid.Nil {
		return models.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}
