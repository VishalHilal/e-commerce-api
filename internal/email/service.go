package email

import (
	"fmt"
	"net/smtp"

	"github.com/VishalHilal/e-commerce-api/internal/models"
)

type EmailService struct {
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		smtpHost: config.SMTPHost,
		smtpPort: config.SMTPPort,
		username: config.Username,
		password: config.Password,
		from:     config.From,
	}
}

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

func (es *EmailService) SendEmail(msg EmailMessage) error {
	auth := smtp.PlainAuth("", es.username, es.password, es.smtpHost)

	smtpAddr := fmt.Sprintf("%s:%d", es.smtpHost, es.smtpPort)

	err := smtp.SendMail(
		smtpAddr,
		auth,
		es.from,
		msg.To,
		[]byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/%s; charset=UTF-8\r\n\r\n%s",
			msg.To[0], msg.Subject, map[bool]string{true: "html", false: "plain"}[msg.IsHTML], msg.Body)),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (es *EmailService) SendWelcomeEmail(user *models.User) error {
	msg := EmailMessage{
		To:      []string{user.Email},
		Subject: "Welcome to Our E-Commerce Store!",
		Body: fmt.Sprintf(`
			<h2>Welcome, %s!</h2>
			<p>Thank you for registering at our e-commerce store. We're excited to have you on board!</p>
			<p>Your account has been successfully created with the email: %s</p>
			<p>You can now:</p>
			<ul>
				<li>Browse our extensive product catalog</li>
				<li>Add items to your cart</li>
				<li>Place orders and track shipments</li>
				<li>Leave product reviews</li>
			</ul>
			<p>If you have any questions, feel free to contact our support team.</p>
			<p>Happy shopping!</p>
			<p>Best regards,<br>The E-Commerce Team</p>
		`, user.FirstName, user.Email),
		IsHTML: true,
	}

	return es.SendEmail(msg)
}

func (es *EmailService) SendOrderConfirmationEmail(user *models.User, order *models.Order) error {
	msg := EmailMessage{
		To:      []string{user.Email},
		Subject: fmt.Sprintf("Order Confirmation - %s", order.OrderNumber),
		Body: fmt.Sprintf(`
			<h2>Order Confirmed!</h2>
			<p>Dear %s,</p>
			<p>Your order <strong>%s</strong> has been successfully placed and is now being processed.</p>
			
			<h3>Order Details:</h3>
			<p><strong>Order Number:</strong> %s</p>
			<p><strong>Total Amount:</strong> $%.2f</p>
			<p><strong>Shipping Address:</strong> %s</p>
			
			<h3>Order Status:</h3>
			<p>Your order is currently: <strong>%s</strong></p>
			
			<p>You will receive another email when your order ships.</p>
			
			<p>Thank you for your purchase!</p>
			<p>Best regards,<br>The E-Commerce Team</p>
		`, user.FirstName, order.OrderNumber, order.OrderNumber, order.TotalAmount, order.ShippingAddress, order.Status),
		IsHTML: true,
	}

	return es.SendEmail(msg)
}

func (es *EmailService) SendShippingConfirmationEmail(user *models.User, order *models.Order) error {
	msg := EmailMessage{
		To:      []string{user.Email},
		Subject: fmt.Sprintf("Your Order Has Shipped - %s", order.OrderNumber),
		Body: fmt.Sprintf(`
			<h2>Good News! Your Order Has Shipped!</h2>
			<p>Dear %s,</p>
			<p>Your order <strong>%s</strong> has been shipped and is on its way to you.</p>
			
			<h3>Shipping Information:</h3>
			<p><strong>Order Number:</strong> %s</p>
			<p><strong>Shipping Address:</strong> %s</p>
			<p><strong>Estimated Delivery:</strong> 3-5 business days</p>
			
			<p>You can track your order using the order number on our website.</p>
			
			<p>Thank you for your patience!</p>
			<p>Best regards,<br>The E-Commerce Team</p>
		`, user.FirstName, order.OrderNumber, order.OrderNumber, order.ShippingAddress),
		IsHTML: true,
	}

	return es.SendEmail(msg)
}

func (es *EmailService) SendPasswordResetEmail(user *models.User, resetToken string) error {
	msg := EmailMessage{
		To:      []string{user.Email},
		Subject: "Password Reset Request",
		Body: fmt.Sprintf(`
			<h2>Password Reset Request</h2>
			<p>Dear %s,</p>
			<p>We received a request to reset your password for your e-commerce account.</p>
			<p>Click the link below to reset your password:</p>
			<p><a href="https://yourstore.com/reset-password?token=%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Reset Password</a></p>
			<p>If you didn't request this password reset, please ignore this email.</p>
			<p>This link will expire in 1 hour for security reasons.</p>
			<p>If you have any questions, contact our support team.</p>
			<p>Best regards,<br>The E-Commerce Team</p>
		`, user.FirstName, resetToken),
		IsHTML: true,
	}

	return es.SendEmail(msg)
}
