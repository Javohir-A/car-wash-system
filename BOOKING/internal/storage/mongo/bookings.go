package mongostorage

import (
	"context"
	"time"

	"booking-service/genproto/bookings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookingsStorage struct {
	collection *mongo.Collection
}

func NewBookingsStorage(coll *mongo.Collection) *BookingsStorage {
	return &BookingsStorage{
		collection: coll,
	}
}

func (s *BookingsStorage) CreateBooking(ctx context.Context, req *bookings.NewBooking) (*bookings.CreateResp, error) {
	objProID, err := primitive.ObjectIDFromHex(req.ProviderId)
	if err != nil {
		return nil, err
	}
	objSerID, err := primitive.ObjectIDFromHex(req.ServiceId)
	if err != nil {
		return nil, err
	}

	booking := bson.M{
		"user_id":        req.UserId,
		"provider_id":    objProID,
		"service_id":     objSerID,
		"status":         req.Status,
		"scheduled_time": req.ScheduledTime,
		"location":       req.Location,
		"total_price":    req.TotalPrice,
		"created_at":     time.Now().Format(time.RFC3339),
		"updated_at":     time.Now().Format(time.RFC3339),
	}

	result, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	return &bookings.CreateResp{
		Id:        result.InsertedID.(primitive.ObjectID).Hex(),
		CreatedAt: booking["created_at"].(string),
	}, nil
}

func (s *BookingsStorage) GetBooking(ctx context.Context, req *bookings.ID) (*bookings.Booking, error) {
	var booking bookings.Booking
	objID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	err = s.collection.FindOne(ctx, filter).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &booking, nil
}

func (s *BookingsStorage) UpdateBooking(ctx context.Context, req *bookings.NewData) (*bookings.UpdateResp, error) {
	objID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	update := bson.M{
		"$set": bson.M{
			"status":         req.Status,
			"scheduled_time": req.ScheduledTime,
			"location":       req.Location,
			"total_price":    req.TotalPrice,
			"updated_at":     time.Now().Format(time.RFC3339),
		},
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &bookings.UpdateResp{
		Id:        req.Id,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *BookingsStorage) CancelBooking(ctx context.Context, req *bookings.ID) (*bookings.Void, error) {
	objID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	_, err = s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &bookings.Void{}, nil
}

func (s *BookingsStorage) ListBookings(ctx context.Context, req *bookings.Pagination) (*bookings.BookingsList, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(req.Limit))
	findOptions.SetSkip(int64((req.Page - 1) * req.Limit))

	cursor, err := s.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookingsList []*bookings.Booking
	for cursor.Next(ctx) {
		var booking bookings.Booking
		if err := cursor.Decode(&booking); err != nil {
			return nil, err
		}
		bookingsList = append(bookingsList, &booking)
	}

	return &bookings.BookingsList{
		Bookings: bookingsList,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}
