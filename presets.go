package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

type QueryPreset struct {
	title   string
	include []string
	exclude []string
}

var PRESETS = map[string]QueryPreset{
	"panama": {
		title:   "Panama",
		include: []string{"panama", "panamá", "tocumen"},
	},
	"cyprus": {
		title:   "Cyprus",
		include: []string{"cyprus", "nicosia", "lefkosia", "limassol", "lemessos", "larnaka", "paphos"},
	},
	"austria": {
		title:   "Austria",
		include: []string{"austria", "österreich", "vienna", "wien", "linz", "salzburg", "graz", "innsbruck", "klagenfurt", "wels", "dornbirn"},
	},
	"armenia": {
		title:   "Armenia",
		include: []string{"armenia", "yerevan", "gyumri", "vanadzor", "vagharshapat", "abovyan", "kapan", "hrazdan", "armavir", "artashat", "ijevan", "gavar", "goris", "dilijan", "stepanakert", "martuni", "sisian", "alaverdi", "stepanavan", "berd"},
	},
	"oman": {
		title:   "Oman",
		include: []string{"oman", "ad+dakhiliyah", "ad+dhahirah", "batinah+north", "batinah+south", "al+buraymi", "al+wusta", "ash+sharqiyah+north", "ash+sharqiyah+south", "dhofar", "muscat", "musandam"},
	},
	"bahrain": {
		title:   "Bahrain",
		include: []string{"bahrain", "manama", "muharraq", "riffa", "hamad+town", "isa+town"},
	},
	"finland": {
		title:   "Finland",
		include: []string{"finland", "suomi", "helsinki", "tampere", "oulu", "espoo", "vantaa", "turku", "rovaniemi", "jyväskylä", "lahti", "kuopio", "pori", "lappeenranta", "vaasa"},
	},
	"sweden": {
		title:   "Sweden",
		include: []string{"sweden", "sverige", "stockholm", "malmö", "uppsala", "göteborg", "gothenburg"},
	},
	"suriname": {
		title:   "Suriname",
		include: []string{"suriname", "paramaribo"},
	},
	"norway": {
		title:   "Norway",
		include: []string{"norway", "norge", "oslo", "bergen", "trondheim", "stavanger", "drammen", "fredrikstad", "kristiansand", "tromsø", "sandnes", "ålesund", "bodø", "skien", "haugesund", "tønsberg", "arendal", "porsgrunn", "hamar", "larvik", "moss", "sandefjord", "halden", "harstad", "lillehammer", "molde", "gjøvik", "mo+i+rana", "steinkjer", "alta", "lommedalen"},
	},
	"germany": {
		title:   "Germany",
		include: []string{"germany", "deutschland", "berlin", "frankfurt", "munich", "münchen", "hamburg", "cologne", "köln"},
	},
	"netherlands": {
		title:   "Netherlands",
		include: []string{"netherlands", "nederland", "amsterdam", "rotterdam", "hague", "utrecht", "holland", "delft"},
	},
	"ukraine": {
		title:   "Ukraine",
		include: []string{"ukraine", "kiev", "kyiv", "kharkiv", "dnipro", "odesa", "donetsk", "zaporizhia"},
	},
	"japan": {
		title:   "Japan",
		include: []string{"japan", "tokyo", "yokohama", "osaka", "nagoya", "sapporo", "kobe", "kyoto", "fukuoka", "kawasaki", "saitama", "hiroshima", "sendai"},
	},
	"russia": {
		title:   "Russia",
		include: []string{"russia", "moscow", "saint+petersburg", "novosibirsk", "yekaterinburg", "nizhny+novgorod", "samara", "omsk", "kazan", "chelyabinsk", "rostov-on-don", "ufa", "volgograd"},
	},
	"estonia": {
		title:   "Estonia",
		include: []string{"estonia", "eesti", "tallinn", "tartu", "narva", "pärnu", "rakvere", "kohtla-järve", "viljandi", "maardu", "sillamäe"},
	},
	"denmark": {
		title:   "Denmark",
		include: []string{"denmark", "danmark", "copenhagen", "aarhus", "odense", "aalborg"},
	},
	"portugal": {
		title:   "Portugal",
		include: []string{"portugal", "lisbon", "lisboa", "braga", "porto", "aveiro", "coimbra", "funchal", "madeira"},
	},
	"france": {
		title:   "France",
		include: []string{"france", "paris", "marseille", "lyon", "toulouse", "nice", "nantes", "strasbourg", "montpellier", "bordeaux", "lille", "rennes", "reims", "rouen", "toulon", "le+havre", "grenoble", "dijon", "le+mans", "brest,france", "tours"},
	},
	"spain": {
		title:   "Spain",
		include: []string{"spain", "españa", "madrid", "barcelona", "valencia", "seville", "sevilla", "zaragoza", "malaga", "murcia", "palma", "bilbao", "alicante", "cordoba"},
	},
	"italy": {
		title:   "Italy",
		include: []string{"italy", "italia", "rome", "roma", "milan", "naples", "napoli", "turin", "torino", "palermo", "genoa", "genova", "bologna", "florence", "firenze", "bari", "catania", "venice", "verona"},
	},
	"uk": {
		title:   "United Kingdom",
		include: []string{"uk", "england", "scotland", "wales", "northern+ireland", "london", "birmingham", "leeds", "glasgow", "sheffield", "bradford", "manchester", "edinburgh", "liverpool", "bristol", "cardiff", "belfast", "leicester", "wakefield", "coventry", "nottingham", "newcastle"},
	},
	"croatia": {
		title:   "Croatia",
		include: []string{"croatia", "hrvatska", "zagreb", "split", "rijeka", "osijek", "zadar", "pula"},
	},
	"worldwide": {
		title:   "Worldwide",
		include: []string{},
	},
	"china": {
		title:   "China",
		include: []string{"china", "中国", "guangzhou", "shanghai", "beijing", "hangzhou"},
	},
	"india": {
		title:   "India",
		include: []string{"india", "mumbai", "delhi", "bangalore", "hyderabad", "ahmedabad", "chennai", "kolkata", "jaipur", "pune", "gurgaon", "noida"},
	},
	"israel": {
		title:   "Israel",
		include: []string{"israel", "tel+aviv", "jerusalem", "beer+sheva", "beersheva", "netanya", "ramat+gan", "haifa", "herzliya", "rishon"},
	},
	"indonesia": {
		title:   "Indonesia",
		include: []string{"indonesia", "jakarta", "surabaya", "bandung", "medan", "bekasi", "semarang", "tangerang", "depok", "makassar", "palembang"},
	},
	"pakistan": {
		title:   "Pakistan",
		include: []string{"pakistan", "karachi", "lahore", "faisalabad", "rawalpindi", "peshawar", "islamabad"},
	},
	"brazil": {
		title:   "Brazil",
		include: []string{"brazil", "brasil", "são+paulo", "brasília", "salvador", "fortaleza", "belém", "belo+horizonte", "manaus", "curitiba", "recife", "rio+de+janeiro", "maceió", "aracaju", "porto+alegre", "florianópolis", "acre", "alagoas", "amapá", "amazonas", "bahia", "ceará", "distrito+federal", "espírito+santo", "goiás", "maranhão", "mato+grosso", "mato+grosso+do+sul", "minas+gerais", "pará", "paraíba", "paraná", "pernambuco", "piauí", "rio+grande+do+norte", "rio+grande+do+sul", "rondônia", "roraima", "santa+catarina", "sergipe", "tocantins"},
	},
	"nigeria": {
		title:   "Nigeria",
		include: []string{"nigeria", "lagos", "kano", "ibadan", "benin+city", "port+harcourt", "jos", "ilorin", "kaduna"},
	},
	"bangladesh": {
		title:   "Bangladesh",
		include: []string{"bangladesh", "dhaka", "chittagong", "khulna", "rajshahi", "barisal", "sylhet", "rangpur", "comilla", "gazipur"},
	},
	"mexico": {
		title:   "Mexico",
		include: []string{"mexico", "mexico+city", "guadalajara", "puebla", "tijuana", "mexicali", "monterrey", "hermosillo", "zapopan", "ciudad+juarez", "chihuahua", "aguascalientes", "mx"},
	},
	"philippines": {
		title:   "Philippines",
		include: []string{"philippines", "pilipinas", "quezon", "manila", "davao", "caloocan", "cebu", "zamboanga", "bohol", "pasig", "bacolod", "makati", "baguio", "cavite"},
	},
	"luxembourg": {
		title:   "Luxembourg",
		include: []string{"luxembourg", "esch-sur-alzette", "differdange", "dudelange", "ettelbruck", "diekirch", "wiltz", "echternach", "rumelange", "grevenmacher", "bertrange", "mamer", "capellen", "strassen", "diekirch"},
	},
	"egypt": {
		title:   "Egypt",
		include: []string{"egypt", "cairo", "alexandria", "giza", "port+said", "suez", "luxor", "el+mahalla", "asyut", "al+mansurah", "tanda"},
		exclude: []string{",+VA", "Virginia", ",+LA", "Louisiana"},
	},
	"ethiopia": {
		title:   "Ethiopia",
		include: []string{"ethiopia", "addis+ababa", "gondar", "adama", "hawassa", "bahir+dar"},
	},
	"vietnam": {
		title:   "Vietnam",
		include: []string{"vietnam", "viet+nam", "ho+chi+minh", "hanoi", "ha+noi", "hai+phong", "da+nang", "can+tho", "bien+hoa", "nha+trang", "vinh"},
	},
	"iran": {
		title:   "Iran",
		include: []string{"iran", "tehran", "mashhad", "isfahan", "esfahan", "karaj", "shiraz", "tabriz", "qom", "ahvaz", "ahwaz", "kermanshah", "urmia", "rasht", "kerman"},
	},
	"congo kinshasa": {
		title:   "Democratic Republic of the Congo",
		include: []string{"congo+kinshasa", "drc", "cod", "kinshasa", "lubumbashi", "bukavu", "kananga", "goma", "mbuji+mayi", "likasi", "kolwezi", "kalemie", "uvira", "matadi", "moba", "kamina", "kabalo", "fungurume"},
	},
	"congo brazzaville": {
		title:   "Republic of the Congo",
		include: []string{"congo+brazza", "cog", "brazzaville", "djambala", "pointe+noire", "sibiti", "owando", "madingou", "loango", "kinkala", "impfondo", "dolisie"},
	},
	"turkey": {
		title:   "Turkey",
		include: []string{"turkey", "turkiye", "istanbul", "ankara", "izmir", "bursa", "adana", "gaziantep", "konya", "antalya", "kayseri", "mersin", "eskisehir", "samsun", "denizli", "malatya"},
	},
	"thailand": {
		title:   "Thailand",
		include: []string{"thailand", "bangkok", "nonthaburi", "nakhon", "phuket", "pattaya", "chiang+mai"},
	},
	"south africa": {
		title:   "South Africa",
		include: []string{"south+africa", "south+africa", "johannesburg", "cape+town", "rsa", "durban", "port+elizabeth", "pretoria", "nelspruit"},
	},
	"myanmar": {
		title:   "Myanmar",
		include: []string{"myanmar", "burma", "yangon", "rangoon", "mandalay", "nay+pyi+taw", "taunggyi", "bago", "mawlamyine"},
	},
	"tanzania": {
		title:   "Tanzania",
		include: []string{"tanzania", "dar+es+salaam", "mwanza", "arusha", "dodoma", "mbeya", "morogoro", "tanga", "kilimanjaro"},
	},
	"south korea": {
		title:   "Republic of Korea",
		include: []string{"south+korea", "ROK", "korea", "seoul", "busan", "incheon", "daegu", "daejeon", "gwangju", "대한민국", "서울", "서울시"},
	},
	"colombia": {
		title:   "Colombia",
		include: []string{"colombia", "bogota", "medellin", "cali", "barranquilla", "cartagena", "cucuta", "bucaramanga", "ibague", "soledad", "pereira", "santa+marta"},
	},
	"kenya": {
		title:   "Kenya",
		include: []string{"kenya", "nairobi", "mombasa", "kisumu", "nakuru", "eldoret", "kisii", "nyeri", "machakos", "Embu"},
	},
	"argentina": {
		title:   "Argentina",
		include: []string{"argentina", "buenos+aires", "cordoba", "rosario", "mendoza", "la+plata", "tucuman", "mar+del+plata", "salta", "resistencia"},
	},
	"algeria": {
		title:   "Algeria",
		include: []string{"algeria", "algiers", "oran", "constantine", "annaba", "blida", "batna", "djelfa", "setif", "sidi+bel+abbes", "biskra", "tiaret", "relizane", "mostaganem", "tlemcen", "chlef", "jijel"},
	},
	"sudan": {
		title:   "Sudan",
		include: []string{"sudan", "khartoum", "omdurman"},
	},
	"poland": {
		title:   "Poland",
		include: []string{"poland", "polska", "warsaw", "krakow", "lodz", "wroclaw", "poznan", "gdansk", "szczecin", "bydgoszcz", "lublin", "katowice", "bialystok"},
	},
	"canada": {
		title:   "Canada",
		include: []string{"canada", "ottawa", "edmonton", "winnipeg", "vancouver", "toronto", "quebec", "montreal", "mississauga", "calgary"},
	},
	"australia": {
		title:   "Australia",
		include: []string{"australia", "sydney", "melbourne", "brisbane", "perth", "adelaide", "canberra", "hobart"},
	},
	"new zealand": {
		title:   "New Zealand",
		include: []string{"new+zealand", "auckland", "wellington", "christchurch", "hamilton", "tauranga", "napier-hastings", "dunedin", "palmerston+north", "nelson", "rotorua", "whangarei", "new+plymouth", "invercargill", "whanganui", "gisborne"},
	},
	"belgium": {
		title:   "Belgium",
		include: []string{"belgium", "antwerp", "ghent", "charleroi", "liege", "brussels", "belgique"},
	},
	"greece": {
		title:   "Greece",
		include: []string{"greece", "Ελλάδα", "athens", "thessaloniki", "patras", "heraklion", "larissa", "volos", "rhodes", "ioannina", "chania", "crete"},
		exclude: []string{"GA"},
	},
	"peru": {
		title:   "Peru",
		include: []string{"peru", "lima", "cusco", "cuzco", "ica", "arequipa", "trujillo", "chiclayo", "huancayo", "piura", "chimbote", "iquitos", "juliaca", "cajamarca"},
	},
	"hungary": {
		title:   "Hungary",
		include: []string{"hungary", "magyarország", "budapest", "szeged", "miskolc"},
	},
	"albania": {
		title:   "Albania",
		include: []string{"albania", "tirana", "durres", "vlore", "elbasan", "shkoder"},
	},
	"uganda": {
		title:   "Uganda",
		include: []string{"uganda", "kampala", "mbarara", "mukono", "jinja", "arua", "gulu", "masaka"},
	},
	"zambia": {
		title:   "Zambia",
		include: []string{"zambia", "lusaka", "kitwe", "ndola"},
	},
	"sri lanka": {
		title:   "Sri Lanka",
		include: []string{"sri+lanka", "balangoda", "ratnapura", "colombo", "moratuwa", "negombo", "galle", "jaffna"},
	},
	"singapore": {
		title:   "Singapore",
		include: []string{"singapore"},
	},
	"latvia": {
		title:   "Latvia",
		include: []string{"latvia", "latvija", "riga", "rīga", "kuldiga", "kuldīga", "ventspils", "liepaja", "liepāja", "daugavpils", "jelgava", "jurmala", "jūrmala"},
	},
	"romania": {
		title:   "Romania",
		include: []string{"romania", "bucharest", "cluj", "iasi", "timisoara", "craiova", "brasov", "sibiu", "constanta", "oradea", "galati", "ploesti", "pitesti", "arad", "bacau"},
	},
	"moldova": {
		title:   "Moldova",
		include: []string{"moldova", "chisinau", "tiraspol", "balti", "bender", "ribnita", "cahul", "ungheni", "soroca", "orhei", "dubasari"},
	},
	"belarus": {
		title:   "Belarus",
		include: []string{"belarus", "minsk", "brest,belarus", "grodno", "gomel", "vitebsk", "mogilev", "slutsk", "borisov", "pinsk", "baranovichi", "bobruisk", "soligorsk"},
	},
	"malta": {
		title:   "Malta",
		include: []string{"malta", "birgu", "bormla", "mdina", "qormi", "senglea", "siġġiewi", "valletta", "zabbar", "zebbuġ", "zejtun"},
	},
	"rwanda": {
		title:   "Rwanda",
		include: []string{"rwanda", "kigali", "butare", "muhanga", "ruhengeri", "gisenyi", "nyarugenge", "huye", "musanze", "rubavu", "rwamagana", "kirehe", "kibungo", "ngoma", "nyagatare", "gicumbi", "nyabihu", "kibuye", "karongi", "rusizi", "nyamasheke", "ruhango", "nyanza", "kamonyi", "kicukiro", "gasabo"},
	},
	"saudi arabia": {
		title:   "Saudi Arabia",
		include: []string{"Saudi", "KSA", "Riyadh", "Mecca", "Jeddah", "Dammam"},
	},
	"morocco": {
		title:   "Morocco",
		include: []string{"morocco", "casablanca", "fez", "tangier", "marrakesh", "salé", "meknes", "rabat", "oujda", "kenitra", "agadir", "tetouan", "temara", "safi", "mohammedia", "khouribga", "el+jadida"},
	},
	"uzbekistan": {
		title:   "Uzbekistan",
		include: []string{"uzbekistan", "tashkent", "namangan", "samarkand", "andijan", "nukus", "bukhara", "qarshi", "fergana"},
	},
	"malaysia": {
		title:   "Malaysia",
		include: []string{"malaysia", "kuala+lumpur", "kajang", "klang", "subang", "penang", "ipoh", "selangor", "melaka", "johor", "sabah", "johor+bahru", "shah+alam", "iskandar+puteri"},
	},
	"afghanistan": {
		title:   "Afghanistan",
		include: []string{"afghanistan", "kabul", "kandahar", "herat", "mazar-e-sharif", "jalalabad", "ghazni", "nangarhar", "khost", "zabul", "helmand", "parwan", "farah", "kunar", "wardak", "baghlan", "kunduz", "takhar", "paktia", "paktika"},
	},
	"venezuela": {
		title:   "Venezuela",
		include: []string{"venezuela", "caracas", "maracaibo", "barquisimeto", "guayana", "maturín", "zulia", "bolivar"},
	},
	"ghana": {
		title:   "Ghana",
		include: []string{"ghana", "accra", "kumasi", "sekondi", "ashaiman", "sunyani", "tamale", "tema"},
	},
	"angola": {
		title:   "Angola",
		include: []string{"angola", "luanda", "huambo", "lobito", "benguela"},
	},
	"nepal": {
		title:   "Nepal",
		include: []string{"nepal", "kathmandu", "pokhara", "lalitpur", "bharatpur", "birgunj", "biratnagar", "janakpur", "ghorahi"},
	},
	"yemen": {
		title:   "Yemen",
		include: []string{"yemen", "sana'a", "taiz", "aden", "mukalla", "ibb"},
	},
	"mozambique": {
		title:   "Mozambique",
		include: []string{"mozambique", "maputo", "matola", "nampula", "beira", "sofala", "chimoio", "tete", "quelimane"},
	},
	"ivory coast": {
		title:   "Ivory Coast",
		include: []string{"ivory", "abidjan", "bouaké", "daloa", "yamoussoukro"},
	},
	"cameroon": {
		title:   "Cameroon",
		include: []string{"cameroon", "Douala", "Yaoundé", "Bafoussam", "Bamenda", "Garoua", "Maroua", "Ngaoundéré", "Kumba", "Nkongsamba", "Buea"},
	},
	"taiwan": {
		title:   "Taiwan",
		include: []string{"taiwan", "Taichung", "Kaohsiung", "Taipei", "Taoyuan", "Tainan", "Hsinchu", "Keelung", "Chiayi", "Changhua"},
	},
	"niger": {
		title:   "Niger",
		include: []string{"niger", "Niamey", "Maradi", "Zinder", "Tahoua", "Agadez", "Arlit", "Birni-N'Konni", "Dosso", "Gaya", "Tessaoua"},
	},
	"burkina faso": {
		title:   "Burkina Faso",
		include: []string{"burkina+faso", "Ouagadougou", "Bobo-Dioulasso", "Koudougou", "Banfora", "Ouahigouya", "Pouytenga", "Kaya", "Tenkodogo", "Fada+N'gourma", "Houndé"},
	},
	"mali": {
		title:   "Mali",
		include: []string{"mali", "bamako", "sikasso", "kalabancoro", "koutiala", "ségou", "kayes", "kati", "mopti", "niono"},
	},
	"malawi": {
		title:   "Malawi",
		include: []string{"malawi", "Lilongwe", "Blantyre", "Mzuzu", "Zomba", "Karonga", "Kasungu", "Mangochi", "Salima", "Liwonde", "Balaka"},
	},
	"chile": {
		title:   "Chile",
		include: []string{"chile", "Santiago", "Valparaíso", "Concepción", "La+Serena", "Antofagasta", "Temuco", "Rancagua", "Talca", "Arica", "Chillán"},
	},
	"kazakhstan": {
		title:   "Kazakhstan",
		include: []string{"kazakhstan", "Almaty", "Shymkent", "Karagandy", "Taraz", "Nur-Sultan", "Pavlodar", "Oskemen", "Semey"},
	},
	"guatemala": {
		title:   "Guatemala",
		include: []string{"Guatemala", "mixco", "villa+nueva", "petapa", "Quetzaltenango"},
	},
	"ecuador": {
		title:   "Ecuador",
		include: []string{"ecuador", "Guayaquil", "Quito", "Cuenca", "Machala"},
	},
	"syria": {
		title:   "Syria",
		include: []string{"syria", "سوريا", "damascus", "hama", "aleppo", "homs", "rif+dimashq", "tartus", "latakia", "idlib", "raqqa", "daraa", "alhasakah", "dierezzor", "quneitra", "alsuwayda"},
	},
	"cambodia": {
		title:   "Cambodia",
		include: []string{"cambodia", "phnom", "battambang", "siem+reap", "kampong"},
	},
	"senegal": {
		title:   "Senegal",
		include: []string{"senegal", "dakar", "touba", "thies", "rufisque", "kaolack", "ziguinchor", "tambacounda", "kaffrine", "diourbel"},
	},
	"chad": {
		title:   "Chad",
		include: []string{"chad", "tchad", "n'djamena", "moundou"},
	},
	"somalia": {
		title:   "Somalia",
		include: []string{"somalia", "mogadishu", "hargeisa", "bosaso", "borama", "garowe", "kismayo"},
	},
	"zimbabwe": {
		title:   "Zimbabwe",
		include: []string{"zimbabwe", "harare", "bulawayo", "mutare", "gweru", "kwekwe"},
	},
	"guinea": {
		title:   "Guinea",
		include: []string{"conakry"},
	},
	"benin": {
		title:   "Benin",
		include: []string{"benin", "cotonou", "porto-novo", "abomey"},
	},
	"haiti": {
		title:   "Haiti",
		include: []string{"haiti", "port-au-prince", "cap-haitien", "carrefour", "delmas", "petion-ville"},
	},
	"cuba": {
		title:   "Cuba",
		include: []string{"cuba", "havana", "santiago+de+cuba", "camaguey", "holguin", "guantanamo", "bayamo"},
	},
	"bolivia": {
		title:   "Bolivia",
		include: []string{"bolivia", "santa+cruz+de+la+sierra", "el+alto", "la+paz", "cochabamba", "oruro", "sucre"},
	},
	"tunisia": {
		title:   "Tunisia",
		include: []string{"tunisia", "tunis", "sfax", "sousse", "kairouan", "ariana", "gabes", "bizerte"},
	},
	"south sudan": {
		title:   "South Sudan",
		include: []string{"south sudan", "juba"},
	},
	"burundi": {
		title:   "Burundi",
		include: []string{"burundi", "bujumbura", "gitega"},
	},
	"dominican republic": {
		title:   "Dominican Republic",
		include: []string{"dominican+republic", "republica+dominicana", "santo+domingo", "la+vega", "macoris"},
	},
	"czech republic": {
		title:   "Czech Republic",
		include: []string{"czech", "czechia", "ceska", "prague", "budejovice", "plzen", "karlovy", "ostrava", "brno"},
	},
	"jordan": {
		title:   "Jordan",
		include: []string{"jordan", "amman", "zarqa", "irbid"},
	},
	"azerbaijan": {
		title:   "Azerbaijan",
		include: []string{"azerbaijan", "baku", "sumqayit", "ganja", "lankaran"},
	},
	"uae": {
		title:   "United Arab Emirates",
		include: []string{"uae", "emirates", "dubai", "abu+dhabi", "sharjah", "al+ain", "ajman"},
	},
	"honduras": {
		title:   "Honduras",
		include: []string{"honduras", "tegucigalpa", "san+pedro+sula", "choloma", "la+ceiba", "el+progreso", "choluteca", "comayagua"},
	},
	"tajikistan": {
		title:   "Tajikistan",
		include: []string{"tajikistan", "dushanbe", "khujand"},
	},
	"papua new guinea": {
		title:   "Papua New Guinea",
		include: []string{"papua+new+guinea", "port+moresby", "lae"},
	},
	"serbia": {
		title:   "Serbia",
		include: []string{"serbia", "belgrade", "novi+sad", "nis", "kragujevac", "subotica", "zrenjanin", "pancevo", "cacak", "novi+pazar", "kraljevo", "smederevo"},
	},
	"switzerland": {
		title:   "Switzerland",
		include: []string{"switzerland", "zurich", "zürich", "geneva", "basel", "lausanne", "bern", "winterthur", "lucerne", "gallen", "lugano", "biel", "thun"},
	},
	"togo": {
		title:   "Togo",
		include: []string{"togo", "lome"},
	},
	"sierra leone": {
		title:   "Sierra Leone",
		include: []string{"sierra+leone", "freetown", "makeni", "koidu"},
	},
	"ireland": {
		title:   "Ireland",
		include: []string{"ireland", "dublin", "cork", "limerick", "galway", "waterford+ireland", "drogheda", "dundalk"},
	},
	"hong kong": {
		title:   "Hong Kong",
		include: []string{"hong+kong", "香港", "kowloon", "九龍"},
	},
	"macau": {
		title:   "Macau",
		include: []string{"macau", "macao"},
	},
	"el salvador": {
		title:   "El Salvador",
		include: []string{"el+salvador"},
	},
	"kyrgyzstan": {
		title:   "Kyrgyzstan",
		include: []string{"kyrgyzstan", "bishkek", "osh", "jalal-abad", "karakol", "tokmok"},
	},
	"nicaragua": {
		title:   "Nicaragua",
		include: []string{"nicaragua", "managua", "matagalpa", "chinandega"},
	},
	"turkmenistan": {
		title:   "Turkmenistan",
		include: []string{"turkmenistan", "turkmenabat"},
	},
	"paraguay": {
		title:   "Paraguay",
		include: []string{"paraguay", "asunción", "asuncion", "ciudad+del+este", "san+lorenzo", "luque", "capiata"},
	},
	"laos": {
		title:   "Laos",
		include: []string{"laos", "vientiane", "pakse"},
	},
	"bulgaria": {
		title:   "Bulgaria",
		include: []string{"bulgaria", "sofia", "plovdiv", "varna", "burgas", "ruse", "stara+zagora", "pleven"},
	},
	"lebanon": {
		title:   "Lebanon",
		include: []string{"lebanon", "beirut", "sidon", "tyre", "tripoli", "byblos", "bekaa", "jounieh", "zahle", "baalbek", "nabatieh", "jbeil", "batroun", "achrafieh", "hamra"},
	},
	"libya": {
		title:   "Libya",
		include: []string{"libya", "tripoli", "benghazi", "misrata", "zliten", "bayda"},
		exclude: []string{"lebanon", "greece", "gr"},
	},
	"slovakia": {
		title:   "Slovakia",
		include: []string{"slovakia", "bratislava", "kosice", "presov", "zilina"},
	},
	"slovenia": {
		title:   "Slovenia",
		include: []string{"slovenia", "slovenija", "ljubljana", "maribor", "celje", "kranj", "koper", "velenje", "novo+mesto", "nova+gorica", "krsko", "krško", "murska+sobota", "postojna", "slovenj+gradec"},
	},
	"lithuania": {
		title:   "Lithuania",
		include: []string{"lithuania", "vilnius", "kaunas", "klaipeda", "siauliai", "panevezys", "alytus"},
	},
	"uruguay": {
		title:   "Uruguay",
		include: []string{"uruguay", "montevideo"},
	},
	"united states": {
		title:   "United States",
		include: []string{",+US", "USA", "United+States", "Alabama", ",+AL", "Alaska", ",+AK", "Arizona", ",+AZ", "Arkansas", ",+AR", "California", ",+CA", "Colorado", ",+CO", "Connecticut", ",+CT", "Delaware", ",+DE", "Florida", ",+FL", "Georgia", ",+GA", "Hawaii", ",+HI", "Idaho", ",+ID", "Illinois", ",+IL", "Indiana", ",+IN", "Iowa", ",+IA", "Kansas", ",+KS", "Kentucky", ",+KY", "Louisiana", ",+LA", "Maine", ",+ME", "Maryland", ",+MD", "Massachusetts", ",+MA", "Michigan", ",+MI", "Minnesota", ",+MN", "Mississippi", ",+MS", "Missouri", ",+MO", "Montana", ",+MT", "Nebraska", ",+NE", "Nevada", ",+NV", "New+Hampshire", ",+NH", "New+Jersey", ",+NJ", "New+Mexico", ",+NM", "New+York", ",+NY", "North+Carolina", ",+NC", "North+Dakota", ",+ND", "Ohio", ",+OH", "Oklahoma", ",+OK", "Oregon", ",+OR", "Pennsylvania", ",+PA", "Rhode+Island", ",+RI", "South+Carolina", ",+SC", "South+Dakota", ",+SD", "Tennessee", ",+TN", "Texas", ",+TX", "Utah", ",+UT", "Vermont", ",+VT", "Virginia", ",+VA", "Washington", ",+WA", "West+Virginia", ",+WV", "Wisconsin", ",+WI", "Wyoming", ",+WY", "Los+Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San+Antonio", "San+Diego", "Dallas", "San+Jose", "Austin", "Jacksonville", "Fort+Worth", "Columbus", "Charlotte", "San+Francisco", "Indianapolis", "Seattle", "Denver", "Boston", "El+Paso", "Nashville", "Detroit", "Portland", "Las+Vegas", "Memphis", "Louisville", "Baltimore"},
	},
	"macedonia": {
		title:   "North Macedonia",
		include: []string{"macedonia", "fyrom", "north+macedonia", "mk", "mkd", "ohd", "skp", "skopje", "bitola", "kumanovo", "prilep", "tetovo", "veles", "shtip", "ohrid", "gostivar", "strumica", "kavadarci", "negotino", "berovo", "kratovo", "struga", "valandovo", "demir+kapija", "demir+hisar", "krusheve", "gevgelija"},
	},
	"palestine": {
		title:   "Palestine",
		include: []string{"palestine", "jerusalem", "gaza", "hebron", "jenin", "nablus", "ramallah", "rafah"},
	},
	"mauritania": {
		title:   "Mauritania",
		include: []string{"mauritania", "mauritanie", "nouakchott", "nouadhibou"},
	},
	"botswana": {
		title:   "Botswana",
		include: []string{"botswana", "gaborone", "francistown"},
	},
	"iraq": {
		title:   "Iraq",
		include: []string{"baghdad", "mosul", "basra", "kirkuk", "erbil", "najaf", "karbala", "sulaymaniya", "al-nasiriya", "al-amarah"},
	},
	"qatar": {
		title:   "Qatar",
		include: []string{"Qatar", "Doha"},
	},
	"the bahamas": {
		title:   "The Bahamas",
		include: []string{"Bahamas"},
	},
	"gabon": {
		title:   "Gabon",
		include: []string{"gabon", "Libreville", "Port-gentil", "Franceville", "Oyem", "Moanda"},
	},
	"georgia": {
		title:   "Georgia",
		include: []string{"Tbilisi", "Batumi", "Kutaisi", "Rustavi", "Zugdidi", "Gori", "Poti", "Telavi", "Akhaltsikhe", "Mtskheta", "Ozurgeti", "Sukhumi", "Samtredia", "Marneuli"},
	},
	"kosovo": {
		title:   "Kosovo",
		include: []string{"kosovo", "kosove", "prishtine"},
	},
	"madagascar": {
		title:   "Madagascar",
		include: []string{"madagascar", "antananarivo", "toamasina", "antsiranana", "mahajanga", "fianarantsoa", "toliara", "antsirabe", "ambositra", "ambatondrazaka", "manakara", "sambava", "morondava", "ambanja", "farafangana", "maintirano", "antsalova", "isoa", "mampikony", "ambatolampy", "ambatofinandrahana", "mandritsara", "marovoay", "moramanga", "vangaindrano", "soaindrana", "ikongo", "tamatave", "diego+suarez", "mananjary", "vohemar", "amparafaravola"},
	},
	"mauritius": {
		title:   "Mauritius",
		include: []string{"mauritius", "port+louis", "curepipe", "quatre+bornes", "vacoas-phoenix", "vacoas", "beau-bassin-rose-hill", "beau+bassin", "rose+hill", "mahebourg", "goodlands", "triolet", "bel+air", "flacq", "souillac", "pamplemousses", "grand+baie", "ebene"},
	},
}

func Preset(name string) QueryPreset {
	return PRESETS[name]
}

func PresetTitle(name string) string {
	title := Preset(name).title
	if title == "" {
		title = strings.Title(name)
	}
	return title
}

func PresetChecksum(name string) string {
	hash := sha256.New()
	io.WriteString(hash, fmt.Sprintf("%+v", Preset(name)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
