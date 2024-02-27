ffmpeg -y \
	-r 30 -t 36.80 -loop 1 -i "virgin-of-the-rocks-paris.jpg" \
	-i "/home/mb/gits/mkshort/.cache/b69efcc8b38d6f4a79aec6bd58d859724a3575f9ff7ddc7cc81e487fc5054c1b.png" \
	-i "/home/mb/gits/mkshort/.cache/81a0e5d6fc40a6b71bee0016dcc8a4d239c43fd819002b9110db5afc97e3d95d.png" \
	-i "/home/mb/gits/mkshort/.cache/7324caa9381e84d365a7efb884691533b87cb69cd5456adbdd27a73f710b5dd1.png" \
	-i "/home/mb/gits/mkshort/.cache/d4499138cd2575f3b62930c2e869847bfe39ac0367c28f8ec4acabd6812fda88.png" \
	-i "/home/mb/gits/mkshort/.cache/c3bac7c1410e9cde28750712bc0296d7ef52e1b322d7d0c5a4cb414d24d73b80.png" \
	-i "/home/mb/gits/mkshort/.cache/448d48033b0c9d127db89e79f77b31577f55002c6547e456603ebb0e1100c9a2.png" \
	-i "/home/mb/gits/mkshort/.cache/b8ea2601b5303af6a9530a4685c9831d41d510eefee2c9ff6b8d679b4f452b2d.png" \
	-i "/home/mb/gits/mkshort/.cache/0293807a8527da63a99cd5c74867e361d3da2752fccaae1790df0ccc81adf371.png" \
	-r 30 -t 17.30 -loop 1 -i "virgin-of-the-rocks-london.jpg" \
	-i "/home/mb/gits/mkshort/.cache/90fd4ba0ef42dc5fb9aae2978e05a1c19a13c1cfb9e68f92ee8590f8b069feb7.png" \
	-i "/home/mb/gits/mkshort/.cache/c537c3515e0dcf60f3d937ec9c4ba81f2ff389a91fba9c26b9dc34b30ad515fc.png" \
	-i "/home/mb/gits/mkshort/.cache/36204cf6c5ffd3977c0409342ed891a570ace5cc14efacca039b4d45a9f8d6a0.png" \
	-i "BMC19T1VivaldiSeasonsSpring.mp3" \
	-filter_complex "
		[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0];
		[9:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1];

		[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1];
		[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2];
		[in0,2] [3:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.60,11.60)' [in0,3];
		[in0,3] [4:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,12.30,15.30)' [in0,4];
		[in0,4] [5:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,16.00,19.00)' [in0,5];
		[in0,5] [6:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,19.70,24.70)' [in0,6];
		[in0,6] [7:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,25.40,28.40)' [in0,7];
		[in0,7] [8:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,31.80,35.80)' [in0,9];
		[in1] [10:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.70)' [in1,1];
		[in1,1] [11:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.10,10.60)' [in1,3];
		[in1,3] [12:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,11.30,14.30)' [in1,4];

		[in0,9] [in1,4] concat=n=2:v=1:a=0:unsafe=1 [v];

		[13:a] afade=type=in:start_time=0:duration=4.00, afade=type=out:start_time=50.10:duration=4.00 [a]
	" \
	-pix_fmt yuv420p -r 30 -movflags faststart -map "[v]" -map "[a]" -c:a aac -shortest "/dev/stdout"
