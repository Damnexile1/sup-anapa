package services

import (
	"fmt"
	"log"
)

type NotificationService struct {
	vkToken string
}

func NewNotificationService(vkToken string) *NotificationService {
	return &NotificationService{vkToken: vkToken}
}

func (s *NotificationService) SendBookingNotification(clientName, clientPhone, slotInfo string) error {
	// TODO: Implement actual VK Bot API integration
	// For now, just log the notification
	message := fmt.Sprintf(
		"Новое бронирование!\nКлиент: %s\nТелефон: %s\nСлот: %s",
		clientName, clientPhone, slotInfo,
	)

	log.Printf("NOTIFICATION: %s", message)

	// Placeholder for future VK API call
	// return s.sendVKMessage(message)

	return nil
}

func (s *NotificationService) sendVKMessage(message string) error {
	// TODO: Implement VK Bot API call
	// https://dev.vk.com/method/messages.send
	return nil
}
