// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
)

type ActivityReport struct {
	ID             string   `json:"_id"`
	Timestamp      int      `json:"Timestamp"`
	Qty            int      `json:"Qty"`
	AvgGain        float64  `json:"AvgGain"`
	TopAGain       *float64 `json:"TopAGain,omitempty"`
	TopBGain       *float64 `json:"TopBGain,omitempty"`
	TopCGain       *float64 `json:"TopCGain,omitempty"`
	FearGreedIndex int      `json:"FearGreedIndex"`
}

type CreateProjectInput struct {
	Title       string    `json:"title"`
	Sop         *bool     `json:"sop,omitempty"`
	Description *string   `json:"description,omitempty"`
	Labels      []*string `json:"labels,omitempty"`
	AssignedTo  *string   `json:"assignedTo,omitempty"`
	DueDate     *string   `json:"dueDate,omitempty"`
	Status      *string   `json:"status,omitempty"`
}

type CreateTaskInput struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Status      *string   `json:"status,omitempty"`
	Labels      []*string `json:"labels,omitempty"`
	AssignedTo  *string   `json:"assignedTo,omitempty"`
	DueDate     *string   `json:"dueDate,omitempty"`
	DeferDate   *string   `json:"deferDate,omitempty"`
	Department  *string   `json:"department,omitempty"`
	ProjectID   *string   `json:"projectId,omitempty"`
	Duration    *int      `json:"duration,omitempty"`
}

type CreateUserInput struct {
	FirstName              string  `json:"firstName"`
	LastName               string  `json:"lastName"`
	Email                  string  `json:"email"`
	Password               string  `json:"password"`
	MobileNumber           *string `json:"mobileNumber,omitempty"`
	Role                   string  `json:"role"`
	InvitedBy              *string `json:"invitedBy,omitempty"`
	PreferredContactMethod *string `json:"preferredContactMethod,omitempty"`
}

type FearAndGreedIndex struct {
	Timestamp           int       `json:"Timestamp"`
	Value               string    `json:"Value"`
	ValueClassification string    `json:"ValueClassification"`
	CreatedAt           time.Time `json:"CreatedAt"`
}

type HistoricKlineData struct {
	Opentime int     `json:"opentime"`
	Coins    []*Ohlc `json:"coins"`
}

type HistoricPrices struct {
	Pair      []*Pair   `json:"Pair,omitempty"`
	Timestamp int       `json:"Timestamp"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type HistoricTickerStats struct {
	Timestamp int            `json:"Timestamp"`
	Stats     []*TickerStats `json:"Stats"`
	CreatedAt time.Time      `json:"CreatedAt"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type MarkAsTestedInput struct {
	BotInstanceName string `json:"BotInstanceName"`
	Tested          bool   `json:"Tested"`
}

type Mean struct {
	Avg   float64 `json:"Avg"`
	Count int     `json:"Count"`
}

type MeanInput struct {
	Avg   float64 `json:"Avg"`
	Count int     `json:"Count"`
}

type Mutation struct {
}

type NewActivityReport struct {
	Timestamp      int      `json:"Timestamp"`
	Qty            int      `json:"Qty"`
	AvgGain        float64  `json:"AvgGain"`
	TopAGain       *float64 `json:"TopAGain,omitempty"`
	TopBGain       *float64 `json:"TopBGain,omitempty"`
	TopCGain       *float64 `json:"TopCGain,omitempty"`
	FearGreedIndex int      `json:"FearGreedIndex"`
}

type NewHistoricKlineDataInput struct {
	Opentime int          `json:"Opentime"`
	Coins    []*OHLCInput `json:"Coins"`
}

type NewHistoricPriceInput struct {
	Pairs     []*PairInput `json:"Pairs"`
	Timestamp int          `json:"Timestamp"`
}

type NewHistoricTickerStatsInput struct {
	Timestamp int                 `json:"Timestamp"`
	Stats     []*TickerStatsInput `json:"Stats"`
}

type NewTradeOutcomeReport struct {
	Timestamp        int      `json:"Timestamp"`
	BotName          string   `json:"BotName"`
	PercentageChange float64  `json:"PercentageChange"`
	Balance          float64  `json:"Balance"`
	Symbol           string   `json:"Symbol"`
	Outcome          string   `json:"Outcome"`
	Fee              *float64 `json:"Fee,omitempty"`
	ElapsedTime      int      `json:"ElapsedTime"`
	Volume           float64  `json:"Volume"`
	FearGreedIndex   int      `json:"FearGreedIndex"`
	MarketStatus     string   `json:"MarketStatus"`
}

type Ohlc struct {
	OpenPrice   string `json:"OpenPrice"`
	HighPrice   string `json:"HighPrice"`
	LowPrice    string `json:"LowPrice"`
	ClosePrice  string `json:"ClosePrice"`
	TradeVolume string `json:"TradeVolume"`
	Symbol      string `json:"Symbol"`
}

type OHLCInput struct {
	OpenPrice   string `json:"OpenPrice"`
	HighPrice   string `json:"HighPrice"`
	LowPrice    string `json:"LowPrice"`
	ClosePrice  string `json:"ClosePrice"`
	TradeVolume string `json:"TradeVolume"`
	Symbol      string `json:"Symbol"`
}

type Pair struct {
	Symbol           string  `json:"Symbol"`
	Price            string  `json:"Price"`
	PercentageChange *string `json:"PercentageChange,omitempty"`
}

type PairInput struct {
	Symbol           string  `json:"Symbol"`
	Price            string  `json:"Price"`
	PercentageChange *string `json:"PercentageChange,omitempty"`
}

type Project struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Sop         bool      `json:"sop"`
	Description *string   `json:"description,omitempty"`
	Labels      []*string `json:"labels,omitempty"`
	AssignedTo  *string   `json:"assignedTo,omitempty"`
	DueDate     *string   `json:"dueDate,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
	Tasks       []*Task   `json:"tasks,omitempty"`
}

type ProjectFilterInput struct {
	Sop *bool `json:"sop,omitempty"`
}

type Query struct {
}

type Strategy struct {
	BotInstanceName      string   `json:"BotInstanceName"`
	TradeDuration        int      `json:"TradeDuration"`
	IncrementsAtr        int      `json:"IncrementsATR"`
	LongSMADuration      int      `json:"LongSMADuration"`
	ShortSMADuration     int      `json:"ShortSMADuration"`
	WINCounter           *int     `json:"WINCounter,omitempty"`
	LOSSCounter          *int     `json:"LOSSCounter,omitempty"`
	TIMEOUTGainCounter   *int     `json:"TIMEOUTGainCounter,omitempty"`
	TIMEOUTLossCounter   *int     `json:"TIMEOUTLossCounter,omitempty"`
	NetGainCounter       *int     `json:"NetGainCounter,omitempty"`
	NetLossCounter       *int     `json:"NetLossCounter,omitempty"`
	AccountBalance       float64  `json:"AccountBalance"`
	MovingAveMomentum    float64  `json:"MovingAveMomentum"`
	TakeProfitPercentage *float64 `json:"TakeProfitPercentage,omitempty"`
	StopLossPercentage   *float64 `json:"StopLossPercentage,omitempty"`
	ATRtollerance        *float64 `json:"ATRtollerance,omitempty"`
	FeesTotal            *float64 `json:"FeesTotal,omitempty"`
	Tested               *bool    `json:"Tested,omitempty"`
	Owner                *string  `json:"Owner,omitempty"`
	CreatedOn            int      `json:"CreatedOn"`
}

type StrategyInput struct {
	BotInstanceName      string   `json:"BotInstanceName"`
	TradeDuration        int      `json:"TradeDuration"`
	IncrementsAtr        int      `json:"IncrementsATR"`
	LongSMADuration      int      `json:"LongSMADuration"`
	ShortSMADuration     int      `json:"ShortSMADuration"`
	WINCounter           *int     `json:"WINCounter,omitempty"`
	LOSSCounter          *int     `json:"LOSSCounter,omitempty"`
	TIMEOUTGainCounter   *int     `json:"TIMEOUTGainCounter,omitempty"`
	TIMEOUTLossCounter   *int     `json:"TIMEOUTLossCounter,omitempty"`
	NetGainCounter       *int     `json:"NetGainCounter,omitempty"`
	NetLossCounter       *int     `json:"NetLossCounter,omitempty"`
	AccountBalance       float64  `json:"AccountBalance"`
	MovingAveMomentum    float64  `json:"MovingAveMomentum"`
	TakeProfitPercentage float64  `json:"TakeProfitPercentage"`
	StopLossPercentage   float64  `json:"StopLossPercentage"`
	ATRtollerance        *float64 `json:"ATRtollerance,omitempty"`
	FeesTotal            *float64 `json:"FeesTotal,omitempty"`
	Tested               *bool    `json:"Tested,omitempty"`
	Owner                string   `json:"Owner"`
	CreatedOn            int      `json:"CreatedOn"`
}

type SymbolStats struct {
	Symbol               string   `json:"Symbol"`
	PositionCounts       []*Mean  `json:"PositionCounts"`
	LiquidityEstimate    *Mean    `json:"LiquidityEstimate,omitempty"`
	MaxLiquidityEstimate *float64 `json:"MaxLiquidityEstimate,omitempty"`
	MinLiquidityEstimate *float64 `json:"MinLiquidityEstimate,omitempty"`
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
	Labels      []*string `json:"labels,omitempty"`
	AssignedTo  *string   `json:"assignedTo,omitempty"`
	DueDate     *string   `json:"dueDate,omitempty"`
	DeferDate   *string   `json:"deferDate,omitempty"`
	Department  *string   `json:"department,omitempty"`
	ProjectID   *string   `json:"projectId,omitempty"`
	Duration    *int      `json:"duration,omitempty"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

type TickerStats struct {
	Symbol            string  `json:"Symbol"`
	PriceChange       string  `json:"PriceChange"`
	PriceChangePct    string  `json:"PriceChangePct"`
	QuoteVolume       string  `json:"QuoteVolume"`
	Volume            string  `json:"Volume"`
	TradeCount        int     `json:"TradeCount"`
	HighPrice         string  `json:"HighPrice"`
	LowPrice          string  `json:"LowPrice"`
	LastPrice         string  `json:"LastPrice"`
	LiquidityEstimate *string `json:"LiquidityEstimate,omitempty"`
}

type TickerStatsInput struct {
	Symbol            string  `json:"Symbol"`
	PriceChange       string  `json:"PriceChange"`
	PriceChangePct    string  `json:"PriceChangePct"`
	QuoteVolume       string  `json:"QuoteVolume"`
	Volume            string  `json:"Volume"`
	TradeCount        int     `json:"TradeCount"`
	HighPrice         string  `json:"HighPrice"`
	LowPrice          string  `json:"LowPrice"`
	LastPrice         string  `json:"LastPrice"`
	LiquidityEstimate *string `json:"LiquidityEstimate,omitempty"`
}

type TradeOutcomeReport struct {
	ID               string   `json:"_id"`
	Timestamp        int      `json:"Timestamp"`
	BotName          string   `json:"BotName"`
	PercentageChange float64  `json:"PercentageChange"`
	Balance          float64  `json:"Balance"`
	Symbol           string   `json:"Symbol"`
	Outcome          string   `json:"Outcome"`
	Fee              *float64 `json:"Fee,omitempty"`
	ElapsedTime      int      `json:"ElapsedTime"`
	Volume           float64  `json:"Volume"`
	FearGreedIndex   int      `json:"FearGreedIndex"`
	MarketStatus     string   `json:"MarketStatus"`
}

type UpdateCountersInput struct {
	BotInstanceName    string   `json:"BotInstanceName"`
	WINCounter         *bool    `json:"WINCounter,omitempty"`
	LOSSCounter        *bool    `json:"LOSSCounter,omitempty"`
	TIMEOUTGainCounter *bool    `json:"TIMEOUTGainCounter,omitempty"`
	TIMEOUTLossCounter *bool    `json:"TIMEOUTLossCounter,omitempty"`
	NetGainCounter     *bool    `json:"NetGainCounter,omitempty"`
	NetLossCounter     *bool    `json:"NetLossCounter,omitempty"`
	AccountBalance     float64  `json:"AccountBalance"`
	FeesTotal          *float64 `json:"FeesTotal,omitempty"`
}

type UpdateProjectInput struct {
	ID          string    `json:"id"`
	Title       *string   `json:"title,omitempty"`
	Sop         *bool     `json:"sop,omitempty"`
	Description *string   `json:"description,omitempty"`
	Labels      []*string `json:"labels,omitempty"`
	AssignedTo  *string   `json:"assignedTo,omitempty"`
	DueDate     *string   `json:"dueDate,omitempty"`
	Status      *string   `json:"status,omitempty"`
}

type UpdateTaskInput struct {
	ID          string    `json:"id"`
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Status      *string   `json:"status,omitempty"`
	Labels      []*string `json:"labels,omitempty"`
	AssignedTo  *string   `json:"assignedTo,omitempty"`
	DueDate     *string   `json:"dueDate,omitempty"`
	DeferDate   *string   `json:"deferDate,omitempty"`
	Department  *string   `json:"department,omitempty"`
	ProjectID   *string   `json:"projectId,omitempty"`
	Duration    *int      `json:"duration,omitempty"`
}

type UpdateUserInput struct {
	ID                     string  `json:"id"`
	FirstName              *string `json:"firstName,omitempty"`
	LastName               *string `json:"lastName,omitempty"`
	Email                  *string `json:"email,omitempty"`
	Password               *string `json:"password,omitempty"`
	MobileNumber           *string `json:"mobileNumber,omitempty"`
	VerifiedEmail          *bool   `json:"verifiedEmail,omitempty"`
	VerifiedMobile         *bool   `json:"verifiedMobile,omitempty"`
	Role                   *string `json:"role,omitempty"`
	IsDeleted              bool    `json:"isDeleted"`
	OpenToTrade            *bool   `json:"openToTrade,omitempty"`
	BinanceAPI             *string `json:"binanceAPI,omitempty"`
	PreferredContactMethod *string `json:"preferredContactMethod,omitempty"`
	Notes                  *string `json:"notes,omitempty"`
	InvitedBy              *string `json:"invitedBy,omitempty"`
	JoinedBallot           *bool   `json:"joinedBallot,omitempty"`
	IsPaidMember           *bool   `json:"isPaidMember,omitempty"`
}

type UpsertFearAndGreedIndexInput struct {
	Timestamp           int    `json:"Timestamp"`
	Value               string `json:"Value"`
	ValueClassification string `json:"ValueClassification"`
}

type UpsertSymbolStatsInput struct {
	Symbol               string       `json:"Symbol"`
	PositionCounts       []*MeanInput `json:"PositionCounts,omitempty"`
	LiquidityEstimate    *MeanInput   `json:"LiquidityEstimate,omitempty"`
	MaxLiquidityEstimate *float64     `json:"MaxLiquidityEstimate,omitempty"`
	MinLiquidityEstimate *float64     `json:"MinLiquidityEstimate,omitempty"`
}

type User struct {
	ID                     string    `json:"id"`
	FirstName              string    `json:"firstName"`
	LastName               string    `json:"lastName"`
	Email                  string    `json:"email"`
	Password               string    `json:"password"`
	MobileNumber           *string   `json:"mobileNumber,omitempty"`
	VerifiedEmail          bool      `json:"verifiedEmail"`
	VerifiedMobile         bool      `json:"verifiedMobile"`
	Role                   string    `json:"role"`
	IsDeleted              bool      `json:"isDeleted"`
	OpenToTrade            bool      `json:"openToTrade"`
	BinanceAPI             *string   `json:"binanceAPI,omitempty"`
	PreferredContactMethod *string   `json:"preferredContactMethod,omitempty"`
	Notes                  *string   `json:"notes,omitempty"`
	InvitedBy              *string   `json:"invitedBy,omitempty"`
	JoinedBallot           bool      `json:"joinedBallot"`
	IsPaidMember           bool      `json:"isPaidMember"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type ContactMethod string

const (
	ContactMethodEmail    ContactMethod = "EMAIL"
	ContactMethodWhatsapp ContactMethod = "WHATSAPP"
)

var AllContactMethod = []ContactMethod{
	ContactMethodEmail,
	ContactMethodWhatsapp,
}

func (e ContactMethod) IsValid() bool {
	switch e {
	case ContactMethodEmail, ContactMethodWhatsapp:
		return true
	}
	return false
}

func (e ContactMethod) String() string {
	return string(e)
}

func (e *ContactMethod) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContactMethod(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContactMethod", str)
	}
	return nil
}

func (e ContactMethod) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *ContactMethod) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e ContactMethod) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

type UserRole string

const (
	UserRoleGuest      UserRole = "GUEST"
	UserRoleInterested UserRole = "INTERESTED"
	UserRoleMember     UserRole = "MEMBER"
	UserRoleAdmin      UserRole = "ADMIN"
)

var AllUserRole = []UserRole{
	UserRoleGuest,
	UserRoleInterested,
	UserRoleMember,
	UserRoleAdmin,
}

func (e UserRole) IsValid() bool {
	switch e {
	case UserRoleGuest, UserRoleInterested, UserRoleMember, UserRoleAdmin:
		return true
	}
	return false
}

func (e UserRole) String() string {
	return string(e)
}

func (e *UserRole) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserRole", str)
	}
	return nil
}

func (e UserRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *UserRole) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e UserRole) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}
