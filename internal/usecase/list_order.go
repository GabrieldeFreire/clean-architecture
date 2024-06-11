package usecase

import (
	"github.com/GabrieldeFreire/clean-architecture/internal/entity"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := c.OrderRepository.GetAllOrders()
	if err != nil {
		return []OrderOutputDTO{}, err
	}
	dtos := make([]OrderOutputDTO, 0, len(orders))
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      float64(order.Price),
			Tax:        float64(order.Tax),
			FinalPrice: float64(order.FinalPrice),
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}
