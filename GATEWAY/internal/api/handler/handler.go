package handler

type MainHandler interface {
	UserManagement() *UserManagementHandler
	Authentication() *AuthHandler
	ProviderManagement() *ProviderManagementHandler
	ServiceManagement() *ServiceManagementHandler
	BookingsManagement() *BookingHandler
	PaymentManagement() *PaymentHandler
}

type MainHandlerImp struct {
	user    *UserManagementHandler
	auth    *AuthHandler
	pro     *ProviderManagementHandler
	ser     *ServiceManagementHandler
	booking *BookingHandler
	payment *PaymentHandler
}

func NewMainHandler(user *UserManagementHandler, auth *AuthHandler,
	pro *ProviderManagementHandler, ser *ServiceManagementHandler,
	bookings *BookingHandler, payment *PaymentHandler) MainHandler {
	return &MainHandlerImp{
		user:    user,
		auth:    auth,
		pro:     pro,
		ser:     ser,
		booking: bookings,
		payment: payment,
	}
}

func (m *MainHandlerImp) UserManagement() *UserManagementHandler {
	return m.user
}

func (m *MainHandlerImp) Authentication() *AuthHandler {
	return m.auth
}

func (m *MainHandlerImp) ProviderManagement() *ProviderManagementHandler {
	return m.pro
}

func (m *MainHandlerImp) ServiceManagement() *ServiceManagementHandler {
	return m.ser
}

func (m *MainHandlerImp) BookingsManagement() *BookingHandler {
	return m.booking
}

func (m *MainHandlerImp) PaymentManagement() *PaymentHandler {
	return m.payment
}
