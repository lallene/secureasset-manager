package asset

type Service struct {
	repo AssetRepository
}

func NewService(repo AssetRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateAsset(asset *Asset) error {
	return s.repo.Create(asset)
}

func (s *Service) GetAssets() ([]Asset, error) {
	return s.repo.FindAll()
}

func (s *Service) GetAssetByID(id string) (*Asset, error) {
	return s.repo.FindByID(id)
}

func (s *Service) UpdateAsset(asset *Asset) error {
	return s.repo.Update(asset)
}

func (s *Service) DeleteAsset(asset *Asset) error {
	return s.repo.Delete(asset)
}
