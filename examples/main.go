package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

const token = "YOUR_TOKEN"

func init() {
	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID
}

func main() {
	rest()
	sandboxRest()
	stream()
}

func stream() {
	logger := log.New(os.Stdout, "[invest-openapi-go-sdk]", log.LstdFlags)

	client, err := sdk.NewStreamingClient(logger, token)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Запускаем цикл обработки входящих событий. Запускаем асинхронно
	// Сюда будут приходить сообщения по подпискам после вызова соответствующих методов
	// SubscribeInstrumentInfo, SubscribeCandle, SubscribeOrderbook
	go func() {
		err = client.RunReadLoop(func(event interface{}) error {
			logger.Printf("Got event %+v", event)
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Подписка на получение событий по инструменту BBG005DXJS36 (TCS)
	err = client.SubscribeInstrumentInfo("BBG005DXJS36", requestID())
	if err != nil {
		log.Fatalln(err)
	}

	// Подписка на получение свечей по инструменту BBG005DXJS36 (TCS)
	err = client.SubscribeCandle("BBG005DXJS36", sdk.CandleInterval5Min, requestID())
	if err != nil {
		log.Fatalln(err)
	}

	// Подписка на получения стакана по инструменту BBG005DXJS36 (TCS)
	err = client.SubscribeOrderbook("BBG005DXJS36", 10, requestID())
	if err != nil {
		log.Fatalln(err)
	}

	// Приложение завершится через 10секунд.
	// Tip: В боевом приложении лучше обрабатывать сигналы завершения и работать в бесконечном цикле
	time.Sleep(10 * time.Second)

	// Отписка от получения событий по инструменту BBG005DXJS36 (TCS)
	err = client.UnsubscribeInstrumentInfo("BBG005DXJS36", requestID())
	if err != nil {
		log.Fatalln(err)
	}

	// Отписка от получения событий по инструменту BBG005DXJS36 (TCS)
	err = client.UnsubscribeCandle("BBG005DXJS36", sdk.CandleInterval5Min, requestID())
	if err != nil {
		log.Fatalln(err)
	}

	// Отписка от получения стакана по инструменту BBG005DXJS36 (TCS)
	err = client.UnsubscribeOrderbook("BBG005DXJS36", 10, requestID())
	if err != nil {
		log.Fatalln(err)
	}
}

func sandboxRest() {
	// Особенность работы в песочнице состоит в том что перед началом работы надо вызвать метод Register
	// и выставить установить(нарисовать) себе активов в портфеле методами SetCurrencyBalance (валютные активы) и
	// SetPositionsBalance (НЕ валютные активы)
	// Обнулить портфель можно методом Clear
	// Все остальные методы rest клиента так же доступны в песочнице
	client := sdk.NewSandboxRestClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Регистрация в песочнице
	err := client.Register(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Рисуем себе 100500 рублей в портфеле песочницы
	err = client.SetCurrencyBalance(ctx, sdk.RUB, 100500)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Рисуем себе 100 акций TCS в портфеле песочницы
	err = client.SetPositionsBalance(ctx, "BBG005DXJS36", 100)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Очищаем состояние портфеля в песочнице
	err = client.Clear(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func rest() {
	client := sdk.NewRestClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение валютных инструментов
	// Например: USD000UTSTOM - USD, EUR_RUB__TOM - EUR
	currencies, err := client.Currencies(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(currencies)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение фондовых инструментов
	// Например: FXMM - Казначейские облигации США, FXGD - золото
	etfs, err := client.ETFs(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(etfs)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение облигационных инструментов
	// Например: SU24019RMFS0 - ОФЗ 24019
	bonds, err := client.Bonds(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bonds)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение акционных инструментов
	// Например: SBUX - Starbucks Corporation
	stocks, err := client.Stocks(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(stocks)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение инструмента по тикеру, возвращает массив инструментов потому что тикер уникален только в рамках одной биржи
	// но может совпадать на разных биржах у разных кампаний
	// Например: https://www.moex.com/ru/issue.aspx?code=FIVE и https://www.nasdaq.com/market-activity/stocks/FIVE
	// В этом примере получить нужную компанию можно проверив поле Currency
	instruments, err := client.SearchInstrumentByTicker(ctx, "TCS")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(instruments)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение инструмента по FIGI(https://en.wikipedia.org/wiki/Financial_Instrument_Global_Identifier)
	// Узнать FIGI нужного инструмента можно методами указанными выше
	// Например: BBG000B9XRY4 - Apple, BBG005DXJS36 - Tinkoff
	instrument, err := client.SearchInstrumentByFIGI(ctx, "BBG005DXJS36")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(instrument)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение списка операций за период по конкретному инструменту(FIGI)
	// Например: ниже запрашиваются операции за последнюю неделю по инструменту NEE
	operations, err := client.Operations(ctx, time.Now().AddDate(0, 0, -7), time.Now(), "BBG000BJSBJ0")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(operations)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение списка НЕ валютных активов портфеля
	positions, err := client.PositionsPortfolio(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(positions)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение списка валютных активов портфеля
	positionCurrencies, err := client.CurrenciesPortfolio(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(positionCurrencies)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение списка валютных и НЕ валютных активов портфеля, метод является совмещеним PositionsPortfolio и CurrenciesPortfolio
	portfolio, err := client.Portfolio(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(portfolio)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение списка выставленных заявок(ордеров)
	orders, err := client.Orders(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(orders)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение свечей(ордеров)
	// Внимание! Действуют ограничения на промежуток и доступный размер свечей за него
	// Интервал свечи и допустимый промежуток запроса:
	// - 1min [1 minute, 1 day]
	// - 2min [2 minutes, 1 day]
	// - 3min [3 minutes, 1 day]
	// - 5min [5 minutes, 1 day]
	// - 10min [10 minutes, 1 day]
	// - 15min [15 minutes, 1 day]
	// - 30min [30 minutes, 1 day]
	// - hour [1 hour, 7 days]
	// - day [1 day, 1 year]
	// - week [7 days, 2 years]
	// - month [1 month, 10 years]
	// Например получение часовых свечей за последние 24 часа по инструменту BBG005DXJS36 (TCS)
	candles, err := client.Candles(ctx, time.Now().AddDate(0, 0, -1), time.Now(), sdk.CandleInterval1Hour, "BBG005DXJS36")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(candles)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение ордербука(он же стакан) по инструменту
	orderbook, err := client.Orderbook(ctx, 10, "BBG005DXJS36")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(orderbook)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выставление лимитной заявки
	// В примере ниже выставляется заявка на покупку ОДНОЙ акции BBG005DXJS36 (TCS) по цене не выше 20$
	placedOrder, err := client.LimitOrder(ctx, "BBG005DXJS36", 1, sdk.BUY, 20)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(placedOrder)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Отмена ранее выставленной заявки.
	// ID заявки возвращается в структуре PlacedLimitOrder в поле ID в запросе выставления заявки client.LimitOrder
	// или в структуре Order в поле ID в запросе получения заявок client.Orders
	err = client.OrderCancel(ctx, "88320371430")
	if err != nil {
		log.Fatalln(err)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Генерируем уникальный ID для запроса
func requestID() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
