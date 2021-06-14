package command

import (
	"math/rand"

	"helm.sh/helm/v3/pkg/time"
)

var lyrics = `Ground Control to Major Tom
Ground Control to Major Tom
Take your protein pills and put your helmet on
(Ten) Ground Control (Nine) to Major Tom (Eight, seven)
(Six) Commencing (Five) countdown, engines on
(Four, three, two)
Check ignition (One) and may God's love (Lift off) be with you

This is Ground Control to Major Tom
You've really made the grade
And the papers want to know whose shirts you wear
Now it's time to leave the capsule if you dare
This is Major Tom to Ground Control
I'm stepping through the door
And I'm floating in a most peculiar way
And the stars look very different today

For here am I sitting in my tin can
Far above the world
Planet Earth is blue
And there's nothing I can do

Though I'm past one hundred thousand miles
I'm feeling very still
And I think my spaceship knows which way to go
Tell my wife I love her very much
She knows
Ground Control to Major Tom
Your circuit's dead, there's something wrong
Can you hear me, Major Tom?
Can you hear me, Major Tom?
Can you hear me, Major Tom?
Can you-

-Here am I floating 'round my tin can
Far above the moon
Planet Earth is blue
And there's nothing I can do`

var bowieArt = "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdhyyyyhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh\nhhhhhhhhhhhhhhhhhhhhhhhhyyhddddhyyyhyyysyyddhhyyyhhhhhhhhhhhhhhhhhhhhh\nhhhhhhhhhhhhhhhhhhhhyyyhmdddhyyssyhddddhyyhhysyyyyyhhhhhhhhhhhhhhhhhhh\nhyhhhhhhhhhyyyyyyyyyyyhmmhyysshhhdmmmmmmmdddhdhhyyyyhhhhhhhhhhhhhhhhhh\nyhhhyyyyyyyyyyyyyyyyhmmmdyyhhyyyhhdNNNmNNNNNmddhyyyyyyyyyhhhhhhhhhhhhh\nyyyyyyyyyyyyyyyyyyyhmNmhysydhdhhdmNNNNNNNNNNmmddddhhyyyyyyhhhhhhhhhhhh\nyyyyyyyyyyyyyyyyyyymNmdhyhdmmNmmNNNNMMNNMNNNmdyhdddhyyyyyyyyhhhhhhhhhh\nyyyyyyyyyyyyyyyyyydNmddddmmNNmddmNNNNNNNNNmhysyhdmmdyyyyyyyyyhhhhhhhhh\nyyyyyyyyyyyyyyyyyhmNmmmNNNNNmhyhdmNNNNNNNmdyssydmNNNhyyyyyyyyyyhhhhhhh\nyyyyyyyyyyyyyyyyydNNmNNNmhsoosssysooyhdhddddhyhhdmNNmyyyyyyyyyyyhhhhhh\nyyyyyyyyyyyyyyyyydNNNNNd+:--::....`.--ss+sydmmNNNNNNmyyyyyyyyyyyhhhhhh\nyyyyyyyyyyyyyyyyymNNNNd+--..:.``````.-/s:/+sdNNMMNMMNhyyyyyyyyyyhhhhhh\nyyyyyyyyyyyyyyyyydNNNNy+:-..:-``    `-//--/+ymNMMMMMMhyyyyyyyyyyyhhhhh\nyyyyyyyyyyyyyyyyyyhNNNyo/:--:/+:---://+///+oydmNMMMMMdyyyyyyyyyyyhhhhh\nyyyyyyyyyyyyyyyyyyymNNyssooooyhho::+hdhysyyhhddNMMMMMhyyyyyyyyyyyyhhhh\nyyyyyyyyyyyyyyyyyyyhNNysyhyys++//..+y//:yhyyyhhdNMMMMdyyyyyyyyyyyyhhhh\nyyyyyyyyyyyyyyyyyyyyhms+//++/-`.:..//-..-::--+shNNNMNdyyyyyyyyyyyyhhhh\nyyyyyyyyyyyyyyyyyyyyyhh+-....`.-:../+:.....-:ohNNmmNhyyyyyyyyyyyhhhhhh\nyyyyyyyyyyyyyyyyyyyyyhdy+/-..../o/:syo-...-/shNNmmNNyyyyyyyyhhhyhhhhhh\nyyyyyyyyyyyyyyyyyyyyyyhhyo/:-../+oyyyo---:/shmNNmNNmyyyyyyyyhhhhhhhhhh\nyyyyyyyyyyyyyyyyyyyyyyhmho/:--....--:--::/oshmNMMNNdyyyyyyyyhhhhhhhhhh\nyyyyyyyyyyyyyyyyyyyyyydddo//::/osssyyys++ossydNMNNNhyyyyyyyyhhhhhhhhhh\nyyyyyyyyyyyyyyyyyyyyyydmNh++//oyysssyyso+syyhNMMNNNhyyyyyyyyyhhhhhhhhh\nyyyyyyyyyyyyyyyyyyyyyyhNNNdso//+osssoooooyhdNNMMNMNhyyyyyyhhhhhhhhhhhh\nyyyyyyyyyyyyyyyyyyyyyyymMNNdy+:---..--:/ohmNNNNNNMNdyyyyyhhhhhhhhhhhhh\nyyyyyyyyyyyyyyyyyyyyyyhNMMNdsyyo/:///+ohmNNNmmNNMMNdyyyhhhhhhhhhhhhhhh\nhhyyyyyyyyyyyyyyosyyyyhNMMNdooydmmmNNNNMMNmdhmNNMMNdyyhhhhhhhhhhhhhhhh\nyyyyyyyyyyooyyys/--/shhNMMNNh++shdmNNNNNNmdydNMMMMNmmddhhhhhhhhhhhhhhh\nyyyyyyyyys:../syho:../ddNMMNho:/syhhhdmmmddmmNMMMNmhs+/smmmmdddhhhhhhh\nyyyyyyyyysyo:-.:sddo.../yNNydh::+yho/+dmNNMMMNdhs/-.-/sdNmhssNNNmdddhh\nhhyyyyhms--/yy:.`-ody/:-./+mNN:--ss+shMNmmmNdo-..-/oydhs/-:/odNNMMMMNd\nhhhhhmMMNm+..:so:-..+hs/-..-+o.--os+/+mNmy+-`.-:ohy/-..-/shm/hssNMNMMM\nhhhhNMMMyMdho.`-oo/-..://:...-..-oh+-:+/.``.:+ys+:.--/oddyso./-sMmMMMM\nhhhmMMMMydhyys/-../o+:.``.--....:+s/.....-:///:.`.:+yhhs/-.-:shhymMMNM\nhhdMMNMNds:./hNho:...::-...--.....```...--.``..-//+o:..-:+ydhNMNhNMNNN\nhmMMm/dNNmhs-`-+yyy/-``---.-:....``.........-//::...-/oyyhsyddmhNdmNMN\ndMMmNhNNddydmy/-.`..::-..---:`.............---...:+hmmo:::odNNmhmdMMMN\nmMMNdyNydMMhMMNyds/-.-:::-.--`.....`..`..-...-:oyyo+:--/symMMMMMNhMMMM\nmNMhddNdMMMddddddmMN+yo/::::-``...`.......-:::::-.-:oymmNNmmmmNMMdmMMN\nmMMdymNdmmmdh+oNhsMMNmmho+/-.``..```...--::::::/shmNMNNNNMMMMMNMMNhNNM"

var bowie = []string{lyrics, bowieArt}

func Bowie() (messages []string) {

	source := rand.NewSource(time.Now().Time.Unix())
	r := rand.New(source)
	index := r.Int() % len(bowie)

	return []string{bowie[index]}
}
