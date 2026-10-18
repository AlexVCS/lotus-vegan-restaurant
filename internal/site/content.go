package site

// The restaurant's facts, in one place so the header, the pages and the
// footer cannot disagree.

// Phone is one pickup line.
type Phone struct{ Display, Tel string }

var phones = []Phone{
	{Display: "727-549-1688", Tel: "tel:+17275491688"},
	{Display: "727-544-0103", Tel: "tel:+17275440103"},
}

// Service is a delivery service.
type Service struct{ Name, Href, Icon string }

var deliveryServices = []Service{
	{Name: "Uber Eats", Icon: "ubereats", Href: "https://www.ubereats.com/store/lotus-vegetarian-restaurant-pinellas-park/h4izgJ7OTvei46tngubVoA"},
	{Name: "Grubhub", Icon: "grubhub", Href: "https://www.grubhub.com/restaurant/lotus-vegan-restaurant-6575-n-pk-blvd-pinellas-park/5557992"},
	{Name: "DoorDash", Icon: "doordash", Href: "https://www.doordash.com/store/lotus-vegan-restaurant-244457/"},
}

const (
	instagramURL  = "https://www.instagram.com/lotusveganpinellas/"
	facebookURL   = "https://www.facebook.com/LotusVeg/"
	directionsURL = "https://www.google.com/maps/search/?api=1&query=Lotus+Vegan+Restaurant+6575+Park+Blvd+N+Pinellas+Park+FL+33781"
)

// Dish is one photo in the "Our food" gallery: the restaurant's own
// Instagram photography.
type Dish struct {
	Name, Note, Alt, Src, Href string
	Width, Height              int
	// Label is the link's accessible name.
	Label string
}

var dishes = []Dish{
	{Name: "Raw Pasta", Note: "Zucchini noodles · No. 21", Src: "/assets/instagram-raw-pasta.webp", Width: 1080, Height: 1080,
		Alt: "Raw zucchini pasta with tomato cashew sauce, fresh tomato and parsley", Href: "https://www.instagram.com/p/CukRMYsJ-Uh/?img_index=1",
		Label: "Raw Pasta — view the post on Instagram (opens in a new tab)"},
	{Name: "Natural Energy Drink", Note: "Watermelon · beet · ginger", Src: "/assets/instagram-energy-drink.webp", Width: 1080, Height: 1080,
		Alt: "Fresh red watermelon, beet, ginger, lime and mint drink", Href: "https://www.instagram.com/lotusveganpinellas/p/DYs8TkDKP1K/",
		Label: "Natural Energy Drink — view the post on Instagram (opens in a new tab)"},
	{Name: "Portobello Noodle", Note: "Ginger sauce · No. 44", Src: "/assets/portobello-noodle.jpg", Width: 1067, Height: 800,
		Alt: "Portobello noodle dish with vegetables in ginger sauce", Href: instagramURL,
		Label: "Portobello Noodle — view Lotus Vegan on Instagram (opens in a new tab)"},
	{Name: "Coffee Frappuccino", Note: "Iced & creamy", Src: "/assets/instagram-coffee-frappuccino.webp", Width: 1080, Height: 1080,
		Alt: "Tall vegan coffee frappuccino against a deep green background", Href: "https://www.instagram.com/lotusveganpinellas/p/C_L7HpjpNYc/",
		Label: "Coffee Frappuccino — view the post on Instagram (opens in a new tab)"},
	{Name: "Garden Hot Pot", Note: "Tofu & vegetables · No. 30", Src: "/assets/instagram-garden-hot-pot.webp", Width: 1080, Height: 1080,
		Alt: "Garden hot pot with tofu broth, fresh vegetables and rice noodles", Href: "https://www.instagram.com/p/Cp3oBuZJeG4/?img_index=1",
		Label: "Garden Hot Pot — view the post on Instagram (opens in a new tab)"},
	{Name: "Chocolate Milk", Note: "Vegan & chilled", Src: "/assets/instagram-chocolate-milk.webp", Width: 1080, Height: 1080,
		Alt: "Iced vegan chocolate milk topped with chocolate foam", Href: "https://www.instagram.com/lotusveganpinellas/p/C_L7HpjpNYc/",
		Label: "Vegan Chocolate Milk — view the post on Instagram (opens in a new tab)"},
	{Name: "Vegan Cheese Sticks", Note: "Marinara · No. 9", Src: "/assets/instagram-vegan-cheese-sticks.webp", Width: 567, Height: 567,
		Alt: "Golden vegan cheese sticks served with marinara sauce", Href: "https://www.instagram.com/p/Cr6TWAJpppp/?img_index=5",
		Label: "Vegan Cheese Sticks — view the post on Instagram (opens in a new tab)"},
	{Name: "Sesame Tofu", Note: "Teriyaki glaze · No. 74", Src: "/assets/instagram-sesame-tofu.webp", Width: 1080, Height: 1080,
		Alt: "Glazed sesame tofu with brown rice, cucumber, lettuce and tomato", Href: "https://www.instagram.com/p/CvyDabRpjv7/?img_index=6",
		Label: "Sesame Tofu — view the post on Instagram (opens in a new tab)"},
	{Name: "Sugarcane Juice", Note: "Bright & naturally sweet", Src: "/assets/instagram-sugarcane-juice.webp", Width: 480, Height: 480,
		Alt: "Bright green fresh sugarcane juice with cut cane", Href: "https://www.instagram.com/lotusveganpinellas/p/DM3TJ0_pfCB/",
		Label: "Fresh Sugarcane Juice — view the post on Instagram (opens in a new tab)"},
}
