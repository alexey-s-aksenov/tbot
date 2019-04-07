package joke

import (
	"strings"
	"testing"

	xhtml "golang.org/x/net/html"
)

func TestExtractData(t *testing.T) {
	const page = `<article class="quote" data-quote="397779">
	<div class="quote__frame">
	  <header class="quote__header">
		<a class="quote__header_permalink" href="/quote/397779">#397779</a>
		<div class="quote__header_date">
		  11.07.2008 в 18:28
		</div>
	  </header>
	  <div class="quote__body">
		fool panda<br>нахуя бля я работаю по такой дерьмовой профессии
		<br>
		<br>Shenter<br>не говори так
		<br>
		<br>Shenter<br>ты же не какой-нибудь дизайнер
		<br>
		<br>fool panda<br>=-O ]:-&gt;<br>
		<br>08.07.2008 14:51:45, Shenter<br>э... Ты дизайнер что ли? :-[		
		<div class="quote__strips" data-debug="1">
          <h3 class="quote__strips_title">
            Комикс по мотивам цитаты
          </h3>
          <ul class="quote__strips_list">
                          <li class="quote__strips_item">
                <a href="/strip/20100104" class="quote__strips_link">
                  <img src="/img/ts/lm5rujqz7uayvieo405422.jpg" class="quote__strips_img">
                </a>
              </li>
                      </ul>
        </div>
	  </div>
	  <footer class="quote__footer">
		<div class="quote__button" role="button" aria-label="Это баян!" tabindex="0" data-vote="[397779, 2, 0]"
			 title="Это баян!">
		  <svg class="quote__dismiss" xmlns="http://www.w3.org/2000/svg" width="50" height="50" fill="none"
			   stroke-width="2" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
			   aria-hidden="true">
			<path d="M21 18H18V32H21"></path>
			<circle cx="22" cy="23" r="1" fill="currentColor" stroke="none"></circle>
			<circle cx="22" cy="27" r="1" fill="currentColor" stroke="none"></circle>
			<path d="M30 17V33"></path>
			<path d="M26 17V33"></path>
			<circle cx="34" cy="23" r="1" fill="currentColor" stroke="none"></circle>
			<circle cx="34" cy="27" r="1" fill="currentColor" stroke="none"></circle>
			<path d="M35 18H38V32H35"></path>
		  </svg>
		</div>
		<div class="quote__vote">
		  <div class="quote__vote_button" role="button" aria-label="Голосовать против этой цитаты" tabindex="0"
			   data-vote="[397779, 1, 0]" data-swipe="1" title="Голосовать против этой цитаты">
			<svg class="quote__vote_icon down" xmlns="http://www.w3.org/2000/svg" width="30" height="30"
				 fill="currentColor" fill-rule="evenodd" clip-rule="evenodd" aria-hidden="true">
			  <path d="M9 14h12v2H9z"></path>
			</svg>
		  </div>
		  <div class="quote__total" data-vote-counter>7728</div>
		  <div class="quote__vote_button" role="button" aria-label="Голосовать за эту цитату" tabindex="0"
			   data-vote="[397779, 0, 0]" data-swipe="0" title="Голосовать за эту цитату">
			<svg class="quote__vote_icon up" xmlns="http://www.w3.org/2000/svg" width="30" height="30"
				 fill="currentColor" fill-rule="evenodd" clip-rule="evenodd" aria-hidden="true">
			  <path d="M16 14v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4z"></path>
			</svg>
		  </div>
		</div>
		<div role="button" class="quote__button" role="button" aria-label="Поделиться" tabindex="0"
			 data-share='{"id": "397779", "url": "quote/397779", "title": "Цитата #397779", "description": "fool panda&lt;br&gt;нахуя бля я работаю по такой дерьмовой профессии&lt;br&gt;&lt;br&gt;Shenter&lt;br&gt;не говори так&lt;br&gt;&lt;br&gt;Shenter&lt;br&gt;ты же не какой-нибудь дизайнер&lt;br&gt;&lt;br&gt;fool panda&lt;br&gt;=-O ]:-&amp;gt;&lt;br&gt;&lt;br&gt;08.07.2008 14:51:45, Shenter&lt;br&gt;э... Ты дизайнер что ли? :-[&lt;br&gt;"}'
			 title="Поделиться цитатой">
		  <svg class="quote__share" xmlns="http://www.w3.org/2000/svg" width="50" height="50" fill="none"
			   stroke-width="2" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
			   aria-hidden="true">
			<path d="M18.5 21l3.5-3.5 3.5 3.5M22 18v10M13 25v7h18v-7"></path>
		  </svg>
		</div>
	  </footer>
	  <div class="quote__overlay">
		<svg class="quote__overlay_icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 30 30"
			 fill="currentColor" fill-rule="evenodd" clip-rule="evenodd" aria-hidden="true">
		  <path class="down" d="M9 14h12v2H9z"></path>
		  <path class="up" d="M16 14v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4z"></path>
		</svg>
	  </div>
	</div>
  </article>`

	doc, _ := xhtml.Parse(strings.NewReader(page))
	var qq = make([]*xhtml.Node, 0)
	//t.Log("testData", nodeToString(doc))
	qq = append(qq, doc)
	actual := extractData(qq)
	var expected = new(quote)
	expected.rating = 7728
	expected.text = `fool panda
нахуя бля я работаю по такой дерьмовой профессии

Shenter
не говори так

Shenter
ты же не какой-нибудь дизайнер

fool panda
=-O ]:-&gt;

08.07.2008 14:51:45, Shenter
э... Ты дизайнер что ли? :-[
`
	expected.text = xhtml.UnescapeString(expected.text)
	if actual == nil {
		t.Fatal("extractQuotes get no resulsts")
	} else {
		if expected.rating != actual[0].rating &&
			expected.text != actual[0].text {
			t.Errorf("Result mismatch\nExpected: r=%v, q=%v\nGot: r=%v, q=%v",
				expected.rating, expected.text,
				actual[0].text, actual[0].rating)
		}
	}
}

type fakeGetWebRequest struct {
}

func (fakeGetWebRequest) FetchBytes(url string) ([]byte, error) {

	return []byte(`<!doctype html>
	<html theme="light" lang="ru" data-controller="random" class="">
	<head>
	  <meta charset="utf-8">
	  <title>Подборка случайных цитат – Цитатник Рунета</title>
	</head>
	<body data-turbolinks-suppress-warning>
	  <div class="columns">
		<main class="columns__main ">
		  <section class="quotes" data-page="">
			<article class="quote" data-quote="405110">
			  <div class="quote__frame">
				<div class="quote__body">
					3<br>3<br>3
				</div>
				<footer class="quote__footer">
	
				  <div class="quote__vote">
					<div class="quote__total" data-vote-counter>33</div>
				  </div>
				</footer>
			  </div>
			</article>
			<article class="quote" data-quote="405110">
			  <div class="quote__frame">
				<div class="quote__body">
					4<br>4<br>4<br>4
				</div>
				<footer class="quote__footer">
	
				  <div class="quote__vote">
					<div class="quote__total" data-vote-counter>4444</div>
				  </div>
				</footer>
			  </div>
			</article>
			<article class="quote" data-quote="418738">
			  <div class="quote__frame">
				<header class="quote__header">
				  <a class="quote__header_permalink" href="/quote/418738">#418738</a>
				  <div class="quote__header_date">
					04.09.2012 в 9:45
				  </div>
				</header>
				<div class="quote__body">
					5<br>5<br>5<br>5<br>5
				</div>
				<footer class="quote__footer">
				  <div class="quote__vote">
					<div class="quote__vote_button" role="button" aria-label="Голосовать против этой цитаты" tabindex="0"
						 data-vote="[418738, 1, 0]" data-swipe="1" title="Голосовать против этой цитаты">
					  <svg class="quote__vote_icon down" xmlns="http://www.w3.org/2000/svg" width="30" height="30"
						   fill="currentColor" fill-rule="evenodd" clip-rule="evenodd" aria-hidden="true">
						<path d="M9 14h12v2H9z"></path>
					  </svg>
					</div>
					<div class="quote__total" data-vote-counter>555</div>
				  </div>
				</footer>
			  </div>
			</article>
		  </section>
	
	
		  <div class="more">
			<a href="/random?4197" class="more__link">
			  Хочу ещё!
			</a>
		  </div>
		</main>
	  </div>
	  <footer class="footer">
		<nav class="footer__nav">
		  <a href="/faq" class="footer__nav_link ">
			О сайте
		  </a>
		  <a href="/webmaster" class="footer__nav_link ">
			Вебмастеру
		  </a>
		  <a href="/rss" class="footer__nav_link">RSS</a>
		  <a href="/rss/comics.xml" class="footer__nav_link">RSS комиксов</a>
		</nav>
	  </footer>
	</body>
	</html>
	`), nil
}

func TestInnerGetJoke(t *testing.T) {
	client := fakeGetWebRequest{}
	expected := `4
4
4
4`
	t.Log("Check on empty hash")
	actual, _ := innerGetJoke(client)
	if actual != expected {
		t.Errorf("Wrong result. Want: %v, Got: %v", expected, actual)
	}
	expected2 := `5
5
5
5
5`
	t.Log("Check with one value in hash")
	actual2, _ := innerGetJoke(client)
	if actual2 != expected2 {
		t.Errorf("Wrong result. Want: %v, Got: %v", expected2, actual2)
	}
	t.Log("Check for hash overwhelming")
	for i := 0; i < 20; i++ {
		actual, _ = innerGetJoke(client)
	}
}
