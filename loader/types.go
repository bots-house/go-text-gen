package loader

type Message struct {
	Key   string
	Value string

	vars []string
}

func (msg *Message) IsTemplate() bool {
	return len(msg.vars) > 0
}

func (msg *Message) Vars() []string {
	return msg.vars
}

func newMessage(k, v string) *Message {
	msg := &Message{
		Key:   k,
		Value: v,
	}

	msg.vars = getTemplateVars(v)

	return msg
}

type Language struct {
	Name     string
	File     string
	Messages []*Message
}

func (lang *Language) Get(key string) *Message {
	for _, msg := range lang.Messages {
		if msg.Key == key {
			return msg
		}
	}
	return nil
}

func (lang *Language) Has(key string) bool {
	for _, msg := range lang.Messages {
		if msg.Key == key {
			return true
		}
	}
	return false
}

func (lang *Language) Keys() []string {
	keys := make([]string, len(lang.Messages))
	for i, msg := range lang.Messages {
		keys[i] = msg.Key
	}
	return keys
}

type Bundle struct {
	Default *Language
	All     []*Language
}

func (bundle *Bundle) Get(lang string) *Language {
	for _, l := range bundle.All {
		if l.Name == lang {
			return l
		}
	}

	return nil
}

func (bundle *Bundle) Languages() []string {
	names := make([]string, len(bundle.All))
	for i, lang := range bundle.All {
		names[i] = lang.Name
	}
	return names
}
