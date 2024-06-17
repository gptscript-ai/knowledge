import pytest
from utils import run_test


@pytest.mark.parametrize("question,answer", [
    ("What is CBA NPAT this year?", "$10,188 million or $10,164 million"),
    pytest.param("On what page does the five-year financial summary start?", "285", marks=pytest.mark.skip(reason="page information is not easy to extract")),
    ("What's the address of CBA in Syndey?", "11 Harbour Street"),
    ("What are the top 3 holders of CommBank PERLS XV Capital Notes?", "BNP, HSBC, Citi"),
    ("How much net profit did New Zealand contribute in 2023?", "1,356, million"),
    ("How much net profit did New Zealand contribute in 2022?", "1,265, million"),
    ("How did H2O.ai help CBA?", "world-leading talent"),
])
def test_CBA_Spreads_dataset(setup_CBA_Spreads_dataset, judge_client, question, answer):
    dataset = setup_CBA_Spreads_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What were Total Liabilities at the end of First Quarter 2023?", "1,816, billion"),
    ("How many branches does TD Bank have in Canada?", "1,060"),
    ("How many Active U.S. banking mobile users does TD Bank have?", "4.8 million"),
])
def test_TD_Bank_dataset(setup_TD_Bank_dataset, judge_client, question, answer):
    dataset = setup_TD_Bank_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What was FY22 total revenue?", "63.1 billion"),
    ("What was the code name of the 4th gen xeon processors?", "Sapphire Rapids"),
    ("What's the world record for overclocking?", "9.008 GHz"),
])
def test_intel_dataset(setup_intel_dataset, judge_client, question, answer):
    dataset = setup_intel_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What are some brands in the Tyson portfolio?", "Hillshire, Jimmy Dean"),
    ("What was the primary driver of volume increase?", "improved, internal, production"),
    ("What was 1H22 net interest expense?", "191 million"),
])
def test_tyson_dataset(setup_tyson_dataset, judge_client, question, answer):
    dataset = setup_tyson_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How many cars did Mercedes-Benz sell in 2022?", "2,040,700"),
    ("How many electrified vehicles did Mercedes-Benz sell in 2022?", "333,500"),
    ("What was the net profit?", "14,809 million"),
])
def test_mercedes_dataset(setup_mercedes_dataset, judge_client, question, answer):
    dataset = setup_mercedes_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What are the total revenues and other income reported by Chevron in 2014?", "211,970 million"),
])
def test_chevron2014_10k_dataset(setup_chevron2014_10k_dataset, judge_client, question, answer):
    dataset = setup_chevron2014_10k_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is the speaker's job?", "Director, Research"),
])
@pytest.mark.skip
def test_imagejon5_dataset(setup_imagejon5_dataset, judge_client, question, answer):
    dataset = setup_imagejon5_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("is the 2nd email starred, yes or no?", "no"),
])
@pytest.mark.skip
def test_imagejoni_dataset(setup_imagejoni_dataset, judge_client, question, answer):
    dataset = setup_imagejoni_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What was the revenue of Brazil?", "15,969, million"),
    ("What was the revenue of Mexico?", "27,229, million"),
    ("How did gross profit change YoY for South America?", "11.0%"),
    ("When was the cybersecurity incident announced?", "April 26"),
    ("Did inflation affect gross profit?", "inflation affected gross profit."),
    ("What country had the largest revenue and how much was it?", "Mexico, 27,229, million"),
])
def test_Femsa_dataset(setup_Femsa_dataset, judge_client, question, answer):
    dataset = setup_Femsa_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Do I need to install CUDA or does Driverless AI ships with CUDA?", "Driverless AI ships with CUDA"),
    ("What's the minimum memory requirements?", "64, GB"),
    ("How do I start Driverless AI in Docker? Please include the docker run command.",
     "docker run, --pid=host, --rm, --shm-size=2g"),
])
def test_DAIInstall_dataset(setup_DAIInstall_dataset, judge_client, question, answer):
    dataset = setup_DAIInstall_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("When should 'PNDG' be used in the price field?", "not available but pending"),
])
def test_esma_dataset(setup_esma_dataset, judge_client, question, answer):
    dataset = setup_esma_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How much of Tengizchevroil does Chevron own?", "50%"),
    ("What was the net income for 2022?", "35,608 million or 35,465 million"),
])
def test_chevron2022_dataset(setup_chevron2022_dataset, judge_client, question, answer):
    dataset = setup_chevron2022_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What drove spending reductions?", "Workforce, reduction, data center"),
    ("How many bolt-on acquisitions have been made since 2021?", "13"),
])
def test_equifax_dataset(setup_equifax_dataset, judge_client, question, answer):
    dataset = setup_equifax_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Who is the CFO from the doc?", "Makarand Padalkar"),
    ("What do Oracles revenues comprise of?", "License fees, Maintenance fees, Consulting fees"),
    ("What was operating profit margin in 2022?", "54%"),
])
def test_oracle_dataset(setup_oracle_dataset, judge_client, question, answer):
    dataset = setup_oracle_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What kind of bond is for investing in states?", "Municipal"),
])
@pytest.mark.skip
def test_imagejon2_dataset(setup_imagejon2_dataset, judge_client, question, answer):
    dataset = setup_imagejon2_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What instrument is the toy bear playing?", "snare"),
])
@pytest.mark.skip
def test_imagejon8_dataset(setup_imagejon8_dataset, judge_client, question, answer):
    dataset = setup_imagejon8_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Find missing data of the sequence: 24 _ 32 33 42", "29"),
])
@pytest.mark.skip
def test_imagejona_dataset(setup_imagejona_dataset, judge_client, question, answer):
    dataset = setup_imagejona_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Extract the text in the image", "Congratulations Kate and Luke on your upcoming arrival"),
    ("Extract the text shown.", "Congratulations Kate and Luke on your upcoming arrival"),
])
@pytest.mark.skip
def test_imagejonk_dataset(setup_imagejonk_dataset, judge_client, question, answer):
    dataset = setup_imagejonk_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Answer question", "January 1, 2013 4:10PM"),
])
@pytest.mark.skip
def test_imagejonq_dataset(setup_imagejonq_dataset, judge_client, question, answer):
    dataset = setup_imagejonq_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How many employees does kaiser permanente have?", "217,000"),
    ("How many lab results were viewed online?", "60.6, million"),
    ("How many members does KP have?", "12.2, million"),
    ("Who's the regional president in Georgia?", "Jim Simpson"),
    ("Who's the CEO at Kaise?", "Greg A. Adams"),
    ("How many nurses work at Kaiser?", "63k"),
    ("How many colorectal cancer screenings happened that year?", "2.2, million"),
])
def test_Kaiser_dataset(setup_Kaiser_dataset, judge_client, question, answer):
    dataset = setup_Kaiser_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How many clients does Bradesco serve?", "77 million"),
    ("Who is the chairman of the board?", "Luiz Carlos Trabuco Cappi"),
    ("What was the number of agreements that include human rights clauses, in 2022?", "22"),
])
def test_Bradesco_dataset(setup_Bradesco_dataset, judge_client, question, answer):
    dataset = setup_Bradesco_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is the Outlook for Eurozone GDP for 2024?", "0.7%"),
    ("What is the Outlook for China GDP for 2023?", "4.9%"),
    ("What is the Outlook for China GDP for 2024?", "4.1%"),
])
def test_Itau_dataset(setup_Itau_dataset, judge_client, question, answer):
    dataset = setup_Itau_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What are the top 3 fast-food restaurants across all age cohorts?", "MCDONALD'S, CHICK-FIL-A, TACO BELL"),
    ("What is the total number of Wendy's customers?", "4,527,294"),
    ("How many baby boomer customers for Subway are there?", "528,785"),
    ("Between which years is a Gen Xer?", "1965-1981"),
    ("Total customers Gen X?", "13,192,015"),
    ("Number of Silent Gen customers for Chipotle?", "16,263"),
    ("Total number of customers for Gen X and Gen Z combined?", "24,038,048"),
])
@pytest.mark.skip
def test_FastFood_dataset(setup_FastFood_dataset, judge_client, question, answer):
    dataset = setup_FastFood_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is whisper?", "speech recognition"),
])
@pytest.mark.skip
def test_DemoDataJon_dataset(setup_DemoDataJon_dataset, judge_client, question, answer):
    dataset = setup_DemoDataJon_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What does rule ID 011 say is the Validation Rule?", "mandatory for all new transaction reports"),
])
@pytest.mark.skip
def test_transxls_dataset(setup_transxls_dataset, judge_client, question, answer):
    dataset = setup_transxls_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How many books did the Adyen team donate to children in-need in San Francisco?", "60"),
])
def test_adyen_dataset(setup_adyen_dataset, judge_client, question, answer):
    dataset = setup_adyen_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What was the third most popular series ever on Netflix?", "Wednesday"),
    ("What was the most popular film in Norway?", "Troll"),
    ("What was the operating margin in 2022?", "18%"),
])
def test_netflix_dataset(setup_netflix_dataset, judge_client, question, answer):
    dataset = setup_netflix_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What percentage of bonds are Municipal Bonds?", "7%"),
])
@pytest.mark.skip
def test_imagejon3_dataset(setup_imagejon3_dataset, judge_client, question, answer):
    dataset = setup_imagejon3_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("When was the revenue highest for newspaper print?", "1999"),
])
@pytest.mark.skip
def test_imagejon7_dataset(setup_imagejon7_dataset, judge_client, question, answer):
    dataset = setup_imagejon7_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What was net income?", "$14,845 million or $14.845 billion"),
    ("How many hours were volunteered, and across how many countries to help confronting society’s challenge?", "115,000, 84"),
    ("How much did Citi finance for affordable housing in the U.S.?", "$6, billion"),
    ("What were total liabilities of Citigroup as of Dec 31 2022?", "$2,214,838 million"),
    ("What were total assets of Citigroup as of Dec 31 2022?", "2,416,676, million"),
    ("On what page are Basel III Revisions?", "49"),
    ("How many employees are at Citi?", "240,000"),
    ("What was the revenue from legacy franchises", "$8.5 billion or $8.472 billion"),
    ("How large is the new stress capital buffer?", "4.0%"),
    ("What were total revenues of Citigroup in 2022?", "$75,338 million or $75.338 billion or $75.3 billion"),
])
def test_CitiAnnual_dataset(setup_CitiAnnual_dataset, judge_client, question, answer):
    dataset = setup_CitiAnnual_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("In which city was Scuderia Ferrari founded and who founded it?", "Modena, Enzo"),
    ("How many cars did Ferrari sell in 2022?", "13,221"),
    ("How many employees did the company have at the end of 2022?", "4,919"),
])
def test_ferrari_dataset(setup_ferrari_dataset, judge_client, question, answer):
    dataset = setup_ferrari_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What was 4th Quarter adjusted net income?", "$20 million"),
    ("How much higher are raw material costs expected to be?", "$300 million"),
    ("Who is the new CFO?", "Christina Zamarro"),
])
def test_goodyear_dataset(setup_goodyear_dataset, judge_client, question, answer):
    dataset = setup_goodyear_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is Jacobs purpose?", "connected, sustainable, world"),
    ("What is Jacobs expected capital expenditure(CAPEX) in 2023?", "125 million"),
    ("What was Critical Mission Solutions revenue in 2022?", "4.4 billion"),
])
def test_jacobs_dataset(setup_jacobs_dataset, judge_client, question, answer):
    dataset = setup_jacobs_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Which tooth in the dental chart is marked with an X?", "21"),
])
@pytest.mark.skip
def test_imagejon6_dataset(setup_imagejon6_dataset, judge_client, question, answer):
    dataset = setup_imagejon6_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What type of foods are in the image?", "fish, carrots"),
    ("What type of foods are shown?", "fish, carrots"),
])
@pytest.mark.skip
def test_imagejon9_dataset(setup_imagejon9_dataset, judge_client, question, answer):
    dataset = setup_imagejon9_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Aidan Gillen acted in how many series?", "2"),
])
@pytest.mark.skip
def test_imagejong_dataset(setup_imagejong_dataset, judge_client, question, answer):
    dataset = setup_imagejong_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is the highest life expectancy at birth of males?", "80.7"),
])
@pytest.mark.skip
def test_imagejonl_dataset(setup_imagejonl_dataset, judge_client, question, answer):
    dataset = setup_imagejonl_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What was goodwill balance?", "$25.2 billion or $25.173 billion"),
    ("What was the available borrowing capacity?", "$209.0, billion"),
    ("What was the value of total foreclosed assets in 2022?", "$137, million"),
    ("How much was the average VaR in 2022?", "$35, million"),
    ("What was long-term debt at the end of 2022?", "$174,870 million"),
    ("What was diluted EPS for 2021?", "$4.95"),
    ("What was diluted EPS for 2022?", "$3.14"),
    ("What was total noninterest income for commercial banking?", "$3,631, million"),
    ("What was total noninterest income for corporate and investment banking?", "$6,509, million"),
    ("What were total nonperforming assets?", "$5,763, million or $5.8 billion"),
])
def test_WellsFargo_dataset(setup_WellsFargo_dataset, judge_client, question, answer):
    dataset = setup_WellsFargo_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How many shares were issued as performance incentive awards in fourth quarter of 2018?", "150 shares"),
    ("What was gross profit in 2017?", "$8,180 million"),
    ("What was total current income tax expense in 2017?", "$1,007 million"),
])
def test_Stryker_dataset(setup_Stryker_dataset, judge_client, question, answer):
    dataset = setup_Stryker_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Who is the CEO of bestbuy?", "Corie Barry"),
    ("How much of the population lives within 10 miles of a Best Buy?", "70%"),
    ("How many totaltech members are there?", "4.6 million"),
])
def test_best_buy_dataset(setup_best_buy_dataset, judge_client, question, answer):
    dataset = setup_best_buy_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How much did DoorDash spend on the gas savings program?", "$40 million"),
    ("Who are the main participants on the call?", "Andy Hargreaves, Prabir Adarkar, Tony Xu"),
])
def test_doordash_dataset(setup_doordash_dataset, judge_client, question, answer):
    dataset = setup_doordash_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    (
    "What is the name of the new suite of componentized and cloud based services that provides banks with highly scalable self-service digital experience capabilities?",
    "Oracle Banking Cloud Services"),
])
def test_ofss_dataset(setup_ofss_dataset, judge_client, question, answer):
    dataset = setup_ofss_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is the RBC value from the report and is that considered abnormal?", "1.8 M/mcL, is considered abnormal"),
])
def test_cbc_sample_report_dataset(setup_cbc_sample_report_dataset, judge_client, question, answer):
    dataset = setup_cbc_sample_report_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What are the total revenues and other income reported by Chevron in 2013?", "228,848 million"),
])
def test_chevron2013_10k_dataset(setup_chevron2013_10k_dataset, judge_client, question, answer):
    dataset = setup_chevron2013_10k_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is the text in the image?", "white cat"),
    ("What is the text shown?", "white cat"),
])
@pytest.mark.skip
def test_imagejon1_dataset(setup_imagejon1_dataset, judge_client, question, answer):
    dataset = setup_imagejon1_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What was the fair amount of paid vacation days in the UK?", "28"),
])
@pytest.mark.skip
def test_imagejonb_dataset(setup_imagejonb_dataset, judge_client, question, answer):
    dataset = setup_imagejonb_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What letter does a keel-shaped cross-section look like?", "V"),
])
@pytest.mark.skip
def test_imagejond_dataset(setup_imagejond_dataset, judge_client, question, answer):
    dataset = setup_imagejond_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("is La Taqueria north of the 24th St Mission Bart station?", "no"),
])
@pytest.mark.skip
def test_imagejonj_dataset(setup_imagejonj_dataset, judge_client, question, answer):
    dataset = setup_imagejonj_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Answer question in the image", "28.01 m/s"),
    ("Answer the question", "28.01 m/s"),
])
@pytest.mark.skip
def test_imagejonp_dataset(setup_imagejonp_dataset, judge_client, question, answer):
    dataset = setup_imagejonp_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How much total assets under management?", "$710, billion"),
    ("How many issuers are in the corporate bond portfolio?", "3,300"),
    ("What was NYLIC's statutory surplus in 2021?", "$24.57, billion"),
    ("What was total surplus (incl. asset valuation reserve)?", "$30.1 billion"),
    ("What percentage is in for Residential Mortgage-Backed Securities?", "6%"),
    ("How large was the dividend payout in 2023?", "$2, billion"),
    ("How large was the general account investment portfolio?", "$317.1, billion"),
    ("Who's America's largest mutual life insurer?", "New York Life"),
    ("When was New York Life insurance founded?", "1845"),
])
def test_NYL_All_dataset(setup_NYL_All_dataset, judge_client, question, answer):
    dataset = setup_NYL_All_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is Label Genie?", "zero-shot labeling"),
    ("Does label genie support audio classification?", "Yes"),
])
@pytest.mark.skip
def test_AudioLabelGenie_dataset(setup_AudioLabelGenie_dataset, judge_client, question, answer):
    dataset = setup_AudioLabelGenie_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is AUM for Franklin by asset class as of September 2022?",
     "Fixed Income,Equity,Alternative,Multi-Asset,Cash Management, $490.9 billion, $392.3 billion, $225.1 billion, $131.5 billion,$57.6 billion"),
    ("How much money was returned to shareholders in 2022?", "$773, million"),
])
def test_franklin_templeton_dataset(setup_franklin_templeton_dataset, judge_client, question, answer):
    dataset = setup_franklin_templeton_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How much was the tax?", "$0.74"),
])
@pytest.mark.skip
def test_imagejon4_dataset(setup_imagejon4_dataset, judge_client, question, answer):
    dataset = setup_imagejon4_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Which is the metro in California that has a good job Outlook?", "Los Angeles"),
])
@pytest.mark.skip
def test_imagejonc_dataset(setup_imagejonc_dataset, judge_client, question, answer):
    dataset = setup_imagejonc_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("If in the food web shown in the diagram, Douglas fir tree needles are absent, which organism would starve?",
     "red tree vole"),
])
@pytest.mark.skip
def test_imagejone_dataset(setup_imagejone_dataset, judge_client, question, answer):
    dataset = setup_imagejone_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What is the name of the tower?", "big ben"),
])
@pytest.mark.skip
def test_imagejonf_dataset(setup_imagejonf_dataset, judge_client, question, answer):
    dataset = setup_imagejonf_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("What was the pack spinner capacity?", "118"),
])
@pytest.mark.skip
def test_imagejonh_dataset(setup_imagejonh_dataset, judge_client, question, answer):
    dataset = setup_imagejonh_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    (
    "Janet Ludlow’s firm requires all its analysts to use a two-stage dividend discount model (DDM) and the capital asset pricing model (CAPM) to value stocks. Using the CAPM and DDM, Ludlow has valued QuickBrush Company at $63 per share. She now must value SmileWhite Corporation. Calculate the required rate of return for SmileWhite by using the information in the following table. A. 14% B. 15% C. 16%",
    "C"),
])
@pytest.mark.skip
def test_imagejonn_dataset(setup_imagejonn_dataset, judge_client, question, answer):
    dataset = setup_imagejonn_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    (
    "Table 11.47 provides a recent survey of the youngest online entrepreneurs whose net worth is estimated at one million dollars or more. Their ages range from 17 to 30. Each cell in the table illustrates the number of entrepreneurs who correspond to the specific age group and their net worth. We want to know whether the ages and net worth independent. \\chi^2 test statistic = ______.  A. 1.56 B. 1.76 C. 1.96 D. 2.06",
    "B"),
])
@pytest.mark.skip
def test_imagejono_dataset(setup_imagejono_dataset, judge_client, question, answer):
    dataset = setup_imagejono_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("Who are the board members?",
     "Christophe Knaub, Yavuz Ölken, Guillaume Herve Marie Xavier Lejeune, Xavier Veyry, Maria Jesus De Arteaga Larru, Nuria Fernandez Paris, Onur Koçkar"),
    ("Compare Axa sigorta's paid claims from 2022 to 2018.", "4,852,940 thousand TL, 2,014,216 thousand TL"),
])
def test_AXA_Sigorta_dataset(setup_AXA_Sigorta_dataset, judge_client, question, answer):
    dataset = setup_AXA_Sigorta_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How much was revenue growth overall?", "30.4%"),
    ("What's the name of the campaign Heineken launched to tackle gender bias?", "Cheers to All Fans"),
    ("What is the leading spirit beer?", "Desperados"),
])
def test_heineken_dataset(setup_heineken_dataset, judge_client, question, answer):
    dataset = setup_heineken_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("How much bonuses was rewarded to frontline associates?", "$580+ million"),
    pytest.param("How many stores are in Florida?", "128", marks=pytest.mark.skip(reason="pdf parser doesn't seem to be able to parse the text")),
    ("What was the adjusted operating margin?", "13.04% or 13.0%"),
])
def test_lowes_dataset(setup_lowes_dataset, judge_client, question, answer):
    dataset = setup_lowes_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("According to the table below, which food is the most likely cause of the outbreak of food poisoning: "
     "A. Cold chicken B. Potato salad C. Egg sandwiches D. Fruit pie and cream", "B"),
])
@pytest.mark.skip
def test_imagejonm_dataset(setup_imagejonm_dataset, judge_client, question, answer):
    dataset = setup_imagejonm_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("what are the 3 main types of images", "Inline, floating and block")
])
def test_docx_1_dataset(setup_docx_1_dataset, judge_client, question, answer):
    dataset = setup_docx_1_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("what are the 3 main types of images", "Inline, floating and block")
])
def test_rtfd_1_dataset(setup_rtfd_1_dataset, judge_client, question, answer):
    dataset = setup_rtfd_1_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("what are the possible tags that can be wrapped to a code block", "pre, code")
])
def test_markdown_dataset(setup_markdown_dataset, judge_client, question, answer):
    dataset = setup_markdown_dataset
    run_test(question, answer, dataset, judge_client)


@pytest.mark.parametrize("question,answer", [
    ("what are the 3 main types of images", "Inline, floating and block")
])
def test_odt_1_dataset(setup_odt_1_dataset, judge_client, question, answer):
    dataset = setup_odt_1_dataset
    run_test(question, answer, dataset, judge_client)
