import time
from selenium import webdriver

driver = webdriver.Firefox()
urlList = []

def getRecipesLinksFromURL(url):
	driver.get(url)
	while True:
		try:
			cookiesDiv = driver.find_element_by_class_name('qc-cmp2-summary-buttons')
			if cookiesDiv is not None:
				cookieButtons = cookiesDiv.find_elements_by_tag_name('button')
				cookieButtons[2].click()
		except Exception as e:
			time.sleep(1)			
		showMoreFound = False
		possibleButtons = driver.find_elements_by_tag_name("button")
		for button in possibleButtons:
			if button.text == "Show more":
				showMoreFound = True
				button.click()
				time.sleep(1)
		if showMoreFound is False:
			elems = driver.find_elements_by_class_name('feed-item')
			for elem in elems:
				urlList.append(elem.get_attribute("href"))
			break

linksWithRecipes = [
	"https://tasty.co/topic/pasta",
	"https://tasty.co/ingredient/chicken",
	"https://tasty.co/ingredient/chocolate",
	"https://tasty.co/ingredient/potato",
	"https://tasty.co/ingredient/ground-beef",
	"https://tasty.co/topic/healthy",
	"https://tasty.co/topic/best-vegetarian",
	"https://tasty.co/topic/low-carb-meals",
	"https://tasty.co/topic/keto",
	"https://tasty.co/topic/vegan",
	"https://tasty.co/topic/breakfast",
	"https://tasty.co/topic/lunch",
	"https://tasty.co/topic/dinner",
	"https://tasty.co/topic/desserts",
	"https://tasty.co/topic/snacks",
	"https://tasty.co/topic/chinese",
	"https://tasty.co/topic/italian",
	"https://tasty.co/topic/japanese"
	"https://tasty.co/topic/mexican"
]

for recipesLink in linksWithRecipes:
	getRecipesLinksFromURL(recipesLink)

driver.close()

urlList = list(set(urlList))

f = open("recipes_links.txt", "w")
for idx in range(len(urlList)):
	if idx == len(urlList) - 1:
		f.write(urlList[idx])
	else:
		f.write(urlList[idx] + "\n")
