# conftest.py
import pytest
import subprocess
import os
from openai import OpenAI


@pytest.fixture(scope="module")
def setup_femsa_dataset():
    datasetName = "Femsa"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/Coca-Cola-FEMSA-Results-1Q23-vf-2.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_wellsfargo_dataset():
    datasetName = "Wellsfargo"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/wellsfargo-2022-annual-report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_citi_annual_dataset():
    datasetName = "CitiAnnual"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/citi-2022-annual-report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_CBA_Spreads_dataset():
    datasetName = "CBA-Spreads"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/2023-Annual-Report-Spreads.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_TD_Bank_dataset():
    datasetName = "TD-Bank"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/2023-td-bank-reports.tar.bz2"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_intel_dataset():
    datasetName = "intel"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(
        ['knowledge', 'ingest', '-d', datasetName, "./data/intel-q4-2022-financial-and-business-report_F.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_tyson_dataset():
    datasetName = "tyson"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(
        ['knowledge', 'ingest', '-d', datasetName, "./data/Tyson-Foods-FINAL-2Q23-Investor-Presentation.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_mercedes_dataset():
    datasetName = "mercedes"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName,
                    "./data/mercedes-benz-annual-report-2022-incl-combined-management-report-mbg-ag.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_chevron2014_10k_dataset():
    datasetName = "chevron2014_10k"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/chevron_2014_10K.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon5_dataset():
    datasetName = "imagejon5"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/jon.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejoni_dataset():
    datasetName = "imagejoni"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/desktop.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_Femsa_dataset():
    datasetName = "Femsa"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/Coca-Cola-FEMSA-Results-1Q23-vf-2.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_DAIInstall_dataset():
    datasetName = "DAIInstall"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/ubuntu.html"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_esma_dataset():
    datasetName = "esma"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(
        ['knowledge', 'ingest', '-d', datasetName, "./data/2016-1452_guidelines_mifid_ii_transaction_reporting.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_chevron2022_dataset():
    datasetName = "chevron2022"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/Chevron-2022-Annual-Report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_equifax_dataset():
    datasetName = "equifax"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName,
                    "./data/equifax-February%2B2023%2BInvestor%2BRelations%2BPresentation.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_oracle_dataset():
    datasetName = "oracle"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/oracle-annual-report-2021-22.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon2_dataset():
    datasetName = "imagejon2"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/ocr2.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon8_dataset():
    datasetName = "imagejon8"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/snare_bear.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejona_dataset():
    datasetName = "imagejona"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/twitter_graph.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonk_dataset():
    datasetName = "imagejonk"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/baby_cake.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonq_dataset():
    datasetName = "imagejonq"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/dates_camps.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_Kaiser_dataset():
    datasetName = "Kaiser"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/kp-annual-report-en-2019.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_Bradesco_dataset():
    datasetName = "Bradesco"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/bradesco-2022-integrated-report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_Itau_dataset():
    datasetName = "Itau"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/Itau_Economic_Prospects_Report-Sep2023.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_FastFood_dataset():
    datasetName = "FastFood"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/fastfood.jpg"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_DemoDataJon_dataset():
    datasetName = "DemoDataJon"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/demo_data_jon.zip"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_transxls_dataset():
    datasetName = "transxls"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName,
                    "./data/esma65-8-2594_annex_1_mifir_transaction_reporting_validation_rules.xlsx"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_adyen_dataset():
    datasetName = "adyen"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/Adyen-Annual-Report-2021.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_netflix_dataset():
    datasetName = "netflix"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/FINAL-Q4-22-Shareholder-Letter.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon3_dataset():
    datasetName = "imagejon3"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/ocr3.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon7_dataset():
    datasetName = "imagejon7"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/revenue.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_CitiAnnual_dataset():
    datasetName = "CitiAnnual"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/citi-2022-annual-report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_ferrari_dataset():
    datasetName = "ferrari"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(
        ['knowledge', 'ingest', '-d', datasetName, "./data/Annual_Report_Ferrari_NV_2022_13.04.2023_Web.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_goodyear_dataset():
    datasetName = "goodyear"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/qtr4_2022_goodyear_investor_letter.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_jacobs_dataset():
    datasetName = "jacobs"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/Jacobs-Investor-Presentation-May-June-2023.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon6_dataset():
    datasetName = "imagejon6"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/dental.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon9_dataset():
    datasetName = "imagejon9"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/fish_and_carrots.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejong_dataset():
    datasetName = "imagejong"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/hbo.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonl_dataset():
    datasetName = "imagejonl"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/chart.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_WellsFargo_dataset():
    datasetName = "WellsFargo"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/wellsfargo-2022-annual-report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_Stryker_dataset():
    datasetName = "Stryker"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/STRYKER_CORPORATION_2018.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_best_buy_dataset():
    datasetName = "best-buy"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/Best-Buy-Investor-Event-March-2022.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_doordash_dataset():
    datasetName = "doordash"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/DASH_Q2-2022-Earnings-Call-Transcript.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_ofss_dataset():
    datasetName = "ofss"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/ofss-annual-report-2022-23.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_cbc_sample_report_dataset():
    datasetName = "cbc_sample_report"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/CBC-sample-report-with-notes_0.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_chevron2013_10k_dataset():
    datasetName = "chevron2013_10k"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/chevron_2013_10K.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon1_dataset():
    datasetName = "imagejon1"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/ocr1.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonb_dataset():
    datasetName = "imagejonb"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/vacation_days.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejond_dataset():
    datasetName = "imagejond"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/leaf_shapes.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonj_dataset():
    datasetName = "imagejonj"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/googlemaps.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonp_dataset():
    datasetName = "imagejonp"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/physics.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_NYL_All_dataset():
    datasetName = "NYL_All"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/2022-nyl-investment-report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_AudioLabelGenie_dataset():
    datasetName = "AudioLabelGenie"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/label-genie-intro-youtube.mp3"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_franklin_templeton_dataset():
    datasetName = "franklin_templeton"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/FRI-2022-Annual-Report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejon4_dataset():
    datasetName = "imagejon4"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/receipt.jpg"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonc_dataset():
    datasetName = "imagejonc"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/jobs.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejone_dataset():
    datasetName = "imagejone"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/red_tree_vole.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonf_dataset():
    datasetName = "imagejonf"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/bigben.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonh_dataset():
    datasetName = "imagejonh"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/displays.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonn_dataset():
    datasetName = "imagejonn"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/janet.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejono_dataset():
    datasetName = "imagejono"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/net_worth.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_AXA_Sigorta_dataset():
    datasetName = "AXA-Sigorta"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/AXA-Sigorta-2022-Annual-Report.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_heineken_dataset():
    datasetName = "heineken"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(
        ['knowledge', 'ingest', '-d', datasetName, "./data/Heineken-NV-Full-Year-press-release-02_15_2023.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_lowes_dataset():
    datasetName = "lowes"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/lowes-2022ar-full-report-4-6-23-final.pdf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_imagejonm_dataset():
    datasetName = "imagejonm"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/food_poisoning.png"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_docx_1_dataset():
    datasetName = "docx"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/docs-demo.docx"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_rtfd_1_dataset():
    datasetName = "rtfd"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/docs-demo.rtf"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_odt_1_dataset():
    datasetName = "odt"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/odt-demo.odt"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="module")
def setup_markdown_dataset():
    datasetName = "markdown"

    subprocess.run(['knowledge', 'create-dataset', datasetName])
    subprocess.run(['knowledge', 'ingest', '-d', datasetName, "./data/markdown-demo.md"])

    yield datasetName

    subprocess.run(['knowledge', 'delete-dataset', datasetName])


@pytest.fixture(scope="session")
def judge_client():
    api_key = os.getenv("OPENAI_API_KEY")
    if api_key == "":
        raise ValueError("env OPENAI_API_KEY is missing")

    base_url = "https://api.openai.com/v1"
    client = OpenAI(base_url=base_url, api_key=api_key)
    yield client
