package minifac

type Resource string

const (
	None    Resource = "none"
	Wood    Resource = "wood"
	Stone   Resource = "stone"
	Coal    Resource = "coal"
	IronOre Resource = "ironore"
	Iron    Resource = "iron"
	Steel   Resource = "steel"
)

func BaseResources() []Resource {
	return []Resource{
		Wood,
		Stone,
		Coal,
		IronOre,
	}
}
