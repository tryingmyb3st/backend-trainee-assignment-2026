package slots_repository

import "backend-assignment-avito/internal/core/domain"

func (r *SlotsRepository) modelsToDomains(models []SlotModel) []domain.Slot {
	slotDomains := make([]domain.Slot, 0, len(models))
	for _, model := range models {
		slotDomains = append(slotDomains, *domain.NewSlot(
			model.ID,
			model.RoomID,
			model.StartTime,
			model.EndTime,
		))
	}
	return slotDomains
}
