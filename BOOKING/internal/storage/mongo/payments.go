package mongostorage

import (
	"booking-service/genproto/payments"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentsStorage struct {
	collection *mongo.Collection
}

func NewPaymentsStorage(col *mongo.Collection) *PaymentsStorage {
	return &PaymentsStorage{
		collection: col,
	}
}

func (s *PaymentsStorage) CreatePayment(ctx context.Context, req *payments.NewPayment) (*payments.CreateResp, error) {
	payment := bson.M{
		"booking_id":     req.BookingId,
		"amount":         req.Amount,
		"status":         req.Status,
		"payment_method": req.PaymentMethod,
		"transaction_id": req.TransactionId,
		"created_at":     time.Now().Format(time.RFC3339),
		"updated_at":     time.Now().Format(time.RFC3339),
	}

	result, err := s.collection.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return &payments.CreateResp{
		Id:        id,
		CreatedAt: payment["created_at"].(string),
	}, nil
}
func (s *PaymentsStorage) ListPayments(ctx context.Context, req *payments.Pagination) (*payments.PaymentsList, error) {
	var paymentsList payments.PaymentsList
	var paymentss []*payments.Payment

	opts := options.Find()
	opts.SetSkip(int64((req.Page - 1) * req.Limit))
	opts.SetLimit(int64(req.Limit))

	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var payment payments.Payment
		if err := cursor.Decode(&payment); err != nil {
			return nil, err
		}

		paymentss = append(paymentss, &payment)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	paymentsList.Payments = paymentss
	paymentsList.Page = req.Page
	paymentsList.Limit = req.Limit

	return &paymentsList, nil
}

func (s *PaymentsStorage) GetPayment(ctx context.Context, req *payments.ID) (*payments.Payment, error) {
	var payment bson.M

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	amount := payment["amount"].(float64)
	return &payments.Payment{
		Id:            payment["_id"].(primitive.ObjectID).Hex(),
		BookingId:     payment["booking_id"].(string),
		Amount:        float32(amount),
		Status:        payment["status"].(string),
		PaymentMethod: payment["payment_method"].(string),
		TransactionId: payment["transaction_id"].(string),
		CreatedAt:     payment["created_at"].(string),
		UpdatedAt:     payment["updated_at"].(string),
	}, nil
}
